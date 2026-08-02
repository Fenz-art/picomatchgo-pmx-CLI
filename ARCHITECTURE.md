# Architecture

## Overview

This repository is a Go port of the core glob-matching engine from picomatch. The implementation is intentionally flat and dependency-free so the package can be used easily from Go applications and tests.

## Core modules

- options.go: defines the public Options configuration and default values.
- types.go: contains shared structs returned by the scanner, parser, and matcher layers.
- scan_impl.go: scans the input pattern to identify prefixes, braces, extglobs, and globstars.
- parse_impl.go: converts glob syntax into regex source and builds the internal parse state.
- matcher_impl.go: wraps the parser output into a public matcher API via MakeRe, CompileRe, and IsMatch.
- utils_impl.go: provides helper functions for regex escaping, POSIX slash normalization, basename handling, and output wrapping.
- constants.go: stores shared constants, character classes, and version information.

## Execution flow

1. A caller invokes Scan, Parse, or IsMatch.
2. Scan extracts structural metadata from the pattern.
3. Parse translates the glob into regex source.
4. CompileRe and MakeRe build a regexp.Regexp object.
5. IsMatch evaluates the final pattern against the target input.

## Design goals

- Preserve behavior from the original picomatch implementation.
- Keep the package dependency-free and idiomatic for Go.
- Maintain small, composable building blocks for parser, scanner, and matcher concerns.
