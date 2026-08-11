# Benchmarks

## Current benchmark targets

The repository includes benchmark-oriented test coverage via the Go test benchmark runner.

Run all benchmarks:

```bash
go test -bench=. -run=^$ ./...
```

## Current results

The following numbers were collected on the current environment during verification:

| Benchmark | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: |
| BenchmarkSimple | 3193 | 0 | 0 |
| BenchmarkGlobstar | 4763 | 0 | 0 |
| BenchmarkBraces | 7258 | 0 | 0 |
| BenchmarkExtglob | 4732 | 0 | 0 |
| BenchmarkPosixClass | 3925 | 0 | 0 |

## Suggested benchmark categories

- simple globs
- globstars
- brace expansion
- extglob matching
- POSIX class matching
- basename matching

## Notes

Benchmark results should be collected on a stable machine and reported alongside the Go version and CPU details.
