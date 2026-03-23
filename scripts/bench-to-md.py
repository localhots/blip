#!/usr/bin/env python3
"""Parse benchstat output and raw benchmark results into clean Markdown tables."""

import re
import sys


def parse_benchstat(text):
    """Parse a benchstat comparison table into structured data."""
    lines = text.strip().splitlines()
    rows = []
    for line in lines:
        # Match lines like: JSON-10  228.0n ± 3%   210.5n ± 1%   -7.66% (p=0.002 n=6)
        m = re.match(
            r"(\S+)\s+"
            r"(\d+\.\d+n)\s*±\s*(\d+%)\s+"
            r"(\d+\.\d+n)\s*±\s*(\d+%)\s+"
            r"([~\-+]\d*\.?\d*%?.*)",
            line,
        )
        if m:
            name = m.group(1).removesuffix("-10")
            rows.append({
                "name": name,
                "base": m.group(2),
                "base_var": m.group(3),
                "curr": m.group(4),
                "curr_var": m.group(5),
                "delta": m.group(6).strip(),
            })
            continue
        # Match lines where only current value exists (no base):
        # JSONNoFields-10   56.25n ± 2%
        m2 = re.match(r"(\S+)\s+(\d+\.\d+n)\s*±\s*(\d+%)\s*$", line)
        if m2:
            name = m2.group(1).removesuffix("-10")
            rows.append({
                "name": name,
                "base": "—",
                "base_var": "",
                "curr": m2.group(2),
                "curr_var": m2.group(3),
                "delta": "—",
            })
            continue
        # Match geomean line
        m3 = re.match(
            r"geomean\s+(\d+\.\d+n)\s+(\d+\.\d+n)\s+([~\-+]\d*\.?\d*%?.*)", line
        )
        if m3:
            rows.append({
                "name": "**geomean**",
                "base": m3.group(1),
                "base_var": "",
                "curr": m3.group(2),
                "curr_var": "",
                "delta": m3.group(3).strip(),
            })
    return rows


def parse_raw_benchmarks(text):
    """Parse raw go test -bench output, dedup by taking first occurrence."""
    lines = text.strip().splitlines()
    seen = {}
    for line in lines:
        m = re.match(
            r"(Benchmark\S+)\s+\d+\s+(\d+\.?\d*)\s+ns/op"
            r"(?:\s+(\d+)\s+B/op)?"
            r"(?:\s+(\d+)\s+allocs/op)?",
            line,
        )
        if m and m.group(1) not in seen:
            name = m.group(1).removesuffix("-10")
            seen[m.group(1)] = {
                "name": name,
                "nsop": float(m.group(2)),
                "bop": int(m.group(3) or 0),
                "allocs": int(m.group(4) or 0),
            }
    return list(seen.values())


def fmt_delta(delta):
    """Format delta string, bolding improvements."""
    delta = delta.strip()
    # Extract just the percentage for cleaner display
    m = re.match(r"([~\-+]\d+\.?\d*%)", delta)
    if m:
        pct = m.group(1)
        if pct.startswith("-"):
            return f"**{pct}**"
        if pct.startswith("~"):
            return "~"
        return pct
    if delta == "—":
        return "—"
    return delta


def render_comparison_table(rows, base_label, curr_label):
    """Render a Markdown table from parsed benchstat rows."""
    lines = []
    lines.append(f"| Benchmark | {base_label} | {curr_label} | Delta |")
    lines.append("|:----------|----------:|----------:|------:|")
    for r in rows:
        base = r["base"]
        if r["base_var"]:
            base += f" ±{r['base_var']}"
        curr = r["curr"]
        if r["curr_var"]:
            curr += f" ±{r['curr_var']}"
        delta = fmt_delta(r["delta"])
        lines.append(f"| {r['name']} | {base} | {curr} | {delta} |")
    return "\n".join(lines)


def render_raw_table(rows):
    """Render a Markdown table of raw benchmark results."""
    lines = []
    lines.append("| Benchmark | ns/op | B/op | allocs/op |")
    lines.append("|:----------|------:|-----:|----------:|")
    for r in rows:
        bop = str(r["bop"]) if r["bop"] > 0 else "0"
        allocs = str(r["allocs"]) if r["allocs"] > 0 else "0"
        lines.append(f"| {r['name']} | {r['nsop']:.1f} | {bop} | {allocs} |")
    return "\n".join(lines)


def main():
    if len(sys.argv) < 4:
        print("Usage: bench-to-md.py <raw_bench_file> <benchstat_vs_main> <benchstat_vs_baseline> <output.md>")
        print("  Pass '-' for any benchstat file to skip that section")
        sys.exit(1)

    raw_file = sys.argv[1]
    vs_main_file = sys.argv[2]
    vs_baseline_file = sys.argv[3]
    output_file = sys.argv[4] if len(sys.argv) > 4 else None

    # Read raw benchmarks
    with open(raw_file) as f:
        raw_text = f.read()
    raw_rows = parse_raw_benchmarks(raw_text)

    # Extract metadata from raw output
    go_version = ""
    cpu = ""
    for line in raw_text.splitlines():
        if line.startswith("cpu:"):
            cpu = line.split(":", 1)[1].strip()
        m = re.match(r"goos: (\w+)", line)

    parts = []
    parts.append("# Benchmark Results\n")
    if cpu:
        parts.append(f"CPU: {cpu}\n")

    # vs main comparison
    if vs_main_file != "-":
        with open(vs_main_file) as f:
            vs_main_text = f.read()
        # Extract only the sec/op section (first table)
        sections = re.split(r"\n\s*│\s+\w+\s+│", vs_main_text)
        rows = parse_benchstat(vs_main_text)
        if rows:
            parts.append("## vs main\n")
            parts.append(render_comparison_table(rows, "main", "current"))
            parts.append("")

    # vs baseline comparison
    if vs_baseline_file != "-":
        with open(vs_baseline_file) as f:
            vs_baseline_text = f.read()
        rows = parse_benchstat(vs_baseline_text)
        if rows:
            parts.append("## vs pre-optimization baseline\n")
            parts.append(render_comparison_table(rows, "baseline", "current"))
            parts.append("")

    # Raw results
    if raw_rows:
        parts.append("## Current Results\n")
        parts.append(render_raw_table(raw_rows))
        parts.append("")

    md = "\n".join(parts)

    if output_file:
        with open(output_file, "w") as f:
            f.write(md)
        print(f"Written to {output_file}")
    else:
        print(md)


if __name__ == "__main__":
    main()
