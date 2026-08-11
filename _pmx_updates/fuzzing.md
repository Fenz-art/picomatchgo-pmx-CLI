# Fuzzing

The repository includes fuzz targets for the core matching and parsing components.

## Run fuzzing locally

```bash
go test -fuzz=FuzzScan -fuzztime=15s ./...
go test -fuzz=FuzzParse -fuzztime=15s ./...
go test -fuzz=FuzzIsMatch -fuzztime=15s ./...
```

## Notes

Fuzzing is useful for surfacing parser edge cases, malformed patterns, and unexpected state transitions.
