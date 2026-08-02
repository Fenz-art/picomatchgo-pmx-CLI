# Benchmarks

## Current benchmark targets

The repository includes benchmark-oriented test coverage via the Go test benchmark runner.

Run all benchmarks:

```bash
go test -bench=. -run=^$ ./...
```

## Suggested benchmark categories

- simple globs
- globstars
- brace expansion
- extglob matching
- POSIX class matching
- basename matching

## Notes

Benchmark results should be collected on a stable machine and reported alongside the Go version and CPU details.
