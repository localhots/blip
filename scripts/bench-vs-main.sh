#!/bin/bash
set -e
set -o pipefail

# Compare current branch benchmarks against stored baselines.
# Baselines are captured once and reused. No branch switching required.
#
# Generates BENCH.md with clean Markdown tables.
#
# Usage: bash scripts/bench-vs-main.sh

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
MAIN_FILE="$SCRIPT_DIR/bench-main.txt"
BASELINE_FILE="$SCRIPT_DIR/bench-optimize-baseline.txt"
BENCH_MD="$PROJECT_DIR/BENCH.md"

command -v benchstat >/dev/null 2>&1 || {
  echo "benchstat not found. Install: go install golang.org/x/perf/cmd/benchstat@latest"
  exit 1
}

TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

echo "Running benchmarks..."
go test -bench=. -benchmem -run='^$' -count=6 2>&1 | tee "$TMP/current.txt"

VS_MAIN_FILE="-"
VS_BASELINE_FILE="-"

# Compare vs main
if [ -f "$MAIN_FILE" ]; then
  echo ""
  echo "=== vs main ==="
  echo ""
  VS_MAIN_FILE="$TMP/vs_main.txt"
  benchstat "$MAIN_FILE" "$TMP/current.txt" 2>&1 | tee "$VS_MAIN_FILE"
fi

# Generate BENCH.md
python3 "$SCRIPT_DIR/bench-to-md.py" \
  "$TMP/current.txt" \
  "$VS_MAIN_FILE" \
  "$VS_BASELINE_FILE" \
  "$BENCH_MD"

echo ""
echo "BENCH.md updated."
