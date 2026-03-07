# Benchmark Results

CPU: Apple M5

## vs main

| Benchmark | main | current | Delta |
|:----------|----------:|----------:|------:|
| JSON | 228.0n ±3% | 199.9n ±0% | **-12.31%** |
| Bare | 141.6n ±2% | 130.7n ±1% | **-7.70%** |
| Optimized | 175.4n ±1% | 163.8n ±2% | **-6.64%** |
| Pretty | 252.4n ±2% | 170.5n ±1% | **-32.44%** |
| PrettySorted | 251.5n ±2% | 174.6n ±1% | **-30.58%** |
| PrettySortedContext | 318.1n ±1% | 254.1n ±2% | **-20.14%** |
| JSONNoFields | — | 49.92n ±1% | — |
| JSONOneField | — | 116.7n ±2% | — |
| JSONTenFields | — | 475.8n ±1% | — |
| JSONContext | — | 292.4n ±1% | — |
| PrettyNoFields | — | 46.44n ±3% | — |
| PrettyTenFields | — | 468.1n ±4% | — |
| Disabled | — | 33.38n ±1% | — |
| ParallelJSON | — | 176.3n ±1% | — |
| ParallelPretty | — | 170.2n ±3% | — |
| **geomean** | 220.3n | 152.7n | **-18.97%** |

## Current Results

| Benchmark | ns/op | B/op | allocs/op |
|:----------|------:|-----:|----------:|
| BenchmarkJSON | 199.6 | 0 | 0 |
| BenchmarkBare | 132.4 | 0 | 0 |
| BenchmarkOptimized | 161.3 | 0 | 0 |
| BenchmarkPretty | 172.9 | 0 | 0 |
| BenchmarkPrettySorted | 175.6 | 0 | 0 |
| BenchmarkPrettySortedContext | 258.4 | 0 | 0 |
| BenchmarkJSONNoFields | 50.3 | 0 | 0 |
| BenchmarkJSONOneField | 117.8 | 0 | 0 |
| BenchmarkJSONTenFields | 476.2 | 616 | 3 |
| BenchmarkJSONContext | 295.7 | 0 | 0 |
| BenchmarkPrettyNoFields | 45.9 | 0 | 0 |
| BenchmarkPrettyTenFields | 466.4 | 616 | 3 |
| BenchmarkDisabled | 33.4 | 0 | 0 |
| BenchmarkParallelJSON | 176.6 | 0 | 0 |
| BenchmarkParallelPretty | 170.5 | 0 | 0 |
