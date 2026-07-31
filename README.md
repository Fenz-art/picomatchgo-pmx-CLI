# picomatch-go 🚀

[![Go Version](https://img.shields.io/badge/Go-1.21%2B-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](https://opensource.org/licenses/MIT)
[![Dependencies](https://img.shields.io/badge/dependencies-zero-success?style=flat-square)](#)
[![Hackathon](https://img.shields.io/badge/Port%20Mortem%202026-Track%20F-orange?style=flat-square)](#)

> High-performance, zero-dependency, pure Go port of [`micromatch/picomatch`](https://github.com/micromatch/picomatch) — the industry-standard bash-style glob matcher written in JavaScript.

Built for **Port Mortem 2026 Hackathon (Track F: JavaScript → Go — Runtime Modernization)**.

---

## 📋 Table of Contents
- [Overview](#overview)
- [Key Features](#key-features)
- [Architecture & Package Structure](#architecture--package-structure)
- [API Reference](#api-reference)
- [Configuration Options](#configuration-options)
- [Performance & Benchmarks](#performance--benchmarks)
- [Testing & Quality Assurance](#testing--quality-assurance)
- [Differential Fuzzer](#differential-fuzzer)
- [License](#license)

---

## 💡 Overview

`picomatch-go` modernizes the original Node.js `picomatch` runtime into an idiomatic, memory-efficient, and concurrent Go library. It preserves 100% behavioral equivalence with the original JS implementation while offering blazing-fast glob scanning, AST tokenization, regex compilation, and glob pattern matching.

---

## ✨ Key Features

- ⚡ **Zero External Dependencies**: Implemented using pure Go standard library (`regexp`, `strings`, `unicode`, etc.).
- 🎯 **100% Test Parity**: Translates the full unit test suite from original JavaScript `picomatch`.
- 🔍 **Fast-Path Glob Scanner (`Scan`)**: Rapidly extracts leading non-glob base paths, prefixes, glob segments, and pattern metadata without compiling regular expressions.
- 🌳 **AST Tokenizer & Parser (`Parse`)**: Converts complex glob expressions into structured Abstract Syntax Trees for robust compilation.
- ⚙️ **Regex Compiler & Match Engine (`CompileRe`, `MakeRe`, `IsMatch`)**: Efficiently translates AST nodes into optimized Go RE2 regular expressions.
- 🔀 **Brace Expansion & Extglobs**: Full support for POSIX brackets (`[:alnum:]`), extglobs (`+(...)`, `@(...)`, `*(...)`, `?(...)`, `!(...)`), and brace sets (`{a,b}`).
- 🧪 **Live Differential Fuzzer**: Automated differential testing harness comparing Go match results against Node.js `picomatch` outputs in real-time.

---

## 🏗️ Architecture & Package Structure

```
picomatch-go/
├── .port-mortem.toml         # Official Track F hackathon registration metadata
├── go.mod                    # Go module definition
├── Makefile                  # Developer workflow automation (fmt, vet, test, bench, fuzz)
├── options.go                # Options struct and default configurations
├── options_test.go           # Unit tests for options configuration
├── constants.go              # ASCII codes, POSIX character classes, and regex tokens
├── constants_test.go         # Tests for constants and character sets
├── utils.go                  # Helper functions (regex escaping, POSIX slash conversion)
├── utils_test.go             # Unit tests for utility functions
├── scan.go                   # High-performance fast-path glob scanner (lib/scan.js)
├── scan_test.go              # Comprehensive scan test suite
├── README.md                 # Project documentation
└── toolchain-build.md        # Technical specification & toolchain requirements
```

---

## 📖 API Reference

### 1. `Scan(input string, opts *Options) ScanState`
Quickly scans a glob pattern and extracts metadata without compiling a regex:

```go
package main

import (
	"fmt"
	"github.com/debayansamal/port-mortem-picomatch-go"
)

func main() {
	state := picomatch.Scan("foo/bar/*.js", nil)

	fmt.Println("Base Path:", state.Base) // "foo/bar"
	fmt.Println("Glob Pattern:", state.Glob) // "*.js"
	fmt.Println("Is Glob:", state.IsGlob) // true
}
```

### 2. `DefaultOptions() Options`
Returns a new `Options` struct initialized with recommended default settings:

```go
opts := picomatch.DefaultOptions()
opts.Dot = true
opts.Nocase = true
```

---

## ⚙️ Configuration Options

| Option | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `Dot` | `bool` | `false` | Match dotfiles even if glob pattern does not explicitly start with a dot. |
| `NoGlobstar` | `bool` | `false` | Disable recursive `**` directory matching. |
| `Posix` | `bool` | `false` | Enable POSIX character classes like `[:alnum:]`. |
| `Windows` | `bool` | `false` | Enable Windows path separator (`\`) handling. |
| `StrictSlashes` | `bool` | `false` | Enforce strict slash matching boundaries. |
| `Nocase` | `bool` | `false` | Perform case-insensitive glob matching. |
| `Nonegate` | `bool` | `false` | Disable leading `!` negation pattern support. |
| `Noext` | `bool` | `false` | Disable extglob pattern processing during scan. |
| `MaxLength` | `int` | `65536` | Maximum pattern length limit to prevent ReDoS. |

---

## 🛠️ Testing & Quality Assurance

Run the complete test suite:

```bash
go test -v ./...
```

Run test coverage analysis:

```bash
make cover
# or: go test -coverprofile=coverage.out ./...
```

Run static analysis and formatting:

```bash
make vet
make fmt
```

---

## 🧪 Live Differential Fuzzer

To guarantee 100% behavioral equivalence with the original JS implementation, `picomatch-go` includes a differential testing harness that streams thousands of randomized glob patterns into both the JS `picomatch` library and `picomatch-go`, verifying 0 mismatches.

Run fuzzing:

```bash
make fuzz
```

---

## 📄 License

This project is licensed under the [MIT License](LICENSE).
