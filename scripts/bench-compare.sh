#!/bin/bash
set -e
set -o pipefail

command -v benchstat >/dev/null 2>&1 || { echo "benchstat not found. Install: go install golang.org/x/perf/cmd/benchstat@latest"; exit 1; }

CURRENT=$(git branch --show-current)
TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

# Stash uncommitted changes so we're comparing clean trees
STASHED=false
if ! git diff --quiet || ! git diff --cached --quiet; then
  git stash push -m "benchcompare temp"
  STASHED=true
fi

# Use short names for benchstat column headers (avoids long temp paths)
BRANCH=$(echo "$CURRENT" | tr '/' '-')
MAIN_FILE="$TMP/main"
BRANCH_FILE="$TMP/$BRANCH"

# Benchmark current branch
echo "Benchmarking $CURRENT..."
go test -bench=. -benchmem -run=Benchmark -count=6 2>&1 | tee "$BRANCH_FILE"

# Switch to main and benchmark
git checkout main
echo "Benchmarking main..."
go test -bench=. -benchmem -run=Benchmark -count=6 2>&1 | tee "$MAIN_FILE"

# Back to original branch
git checkout "$CURRENT"

if $STASHED; then
  git stash pop || { echo "warning: stash pop failed (conflicts?). Resolve with: git stash drop"; exit 1; }
fi

# Compare (main = baseline, current branch = new)
# Run from TMP so benchstat shows "main" and "optimize" instead of full paths
(cd "$TMP" && benchstat main "$BRANCH")
