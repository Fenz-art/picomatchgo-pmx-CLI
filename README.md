# picomatch-go

A high-performance, zero-dependency pure Go port of [`micromatch/picomatch`](https://github.com/micromatch/picomatch).

Built for **Port Mortem 2026 Hackathon (Track F: JavaScript → Go Runtime Modernization)**.

## Project Goals
- Zero external dependencies (Go standard library only)
- 100% test parity with original JS unit tests
- High-performance glob matching & AST compilation
- Live differential fuzzer against original JS library

## Single-Command Build & Test
```bash
go test ./...
```
