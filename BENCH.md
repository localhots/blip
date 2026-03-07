# Benchmark Results

CPU: Apple M5

## Current Results (Mar 2026)

| Benchmark | ns/op | B/op | allocs/op |
|:----------|------:|-----:|----------:|
| BenchmarkJSON | 211.0 | 0 | 0 |
| BenchmarkBare | 129.0 | 0 | 0 |
| BenchmarkOptimized | 162.0 | 0 | 0 |
| BenchmarkPretty | 170.3 | 0 | 0 |
| BenchmarkPrettySorted | 174.9 | 0 | 0 |
| BenchmarkPrettySortedContext | 259.3 | 0 | 0 |
| BenchmarkJSONNoFields | 55.2 | 0 | 0 |
| BenchmarkJSONOneField | 120.4 | 0 | 0 |
| BenchmarkJSONTenFields | 468.9 | 616 | 3 |
| BenchmarkJSONContext | 290.3 | 0 | 0 |
| BenchmarkPrettyNoFields | 47.7 | 0 | 0 |
| BenchmarkPrettyTenFields | 463.0 | 616 | 3 |
| BenchmarkDisabled | 33.6 | 0 | 0 |
| BenchmarkParallelJSON | 184.5 | 0 | 0 |
| BenchmarkParallelPretty | 177.4 | 0 | 0 |

## vs main (Sep 2025)

| Benchmark | main | current | Delta |
|:----------|----------:|----------:|------:|
| JSON | 228.0n ±3% | 216.1n ±8% | **-5.20%** |
| Bare | 141.6n ±2% | 129.0n ±1% | **-8.87%** |
| Optimized | 175.4n ±1% | 159.9n ±1% | **-8.81%** |
| Pretty | 252.4n ±2% | 171.0n ±0% | **-32.26%** |
| PrettySorted | 251.5n ±2% | 175.0n ±1% | **-30.42%** |
| PrettySortedContext | 318.1n ±1% | 247.8n ±5% | **-22.12%** |
| JSONNoFields | — | 55.41n ±2% | — |
| JSONOneField | — | 120.1n ±0% | — |
| JSONTenFields | — | 468.4n ±1% | — |
| JSONContext | — | 288.5n ±1% | — |
| PrettyNoFields | — | 47.74n ±0% | — |
| PrettyTenFields | — | 463.5n ±0% | — |
| Disabled | — | 33.48n ±1% | — |
| ParallelJSON | — | 176.6n ±4% | — |
| ParallelPretty | — | 169.3n ±5% | — |
| **geomean** | 220.3n | 154.1n | **-18.68%** |
