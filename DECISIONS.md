# DECISIONS.md — Design & Porting Decisions

This document records the key design decisions made while porting [`micromatch/picomatch`](https://github.com/micromatch/picomatch) from JavaScript to Go.

---

## 1. Flat Package Layout vs Nested `internal/` Packages

**Decision**: Use a single flat `picomatch` package instead of splitting into `internal/scanner`, `internal/parser`, etc.

**Rationale**: The original JS library exports everything from a single module (`picomatch`). A flat Go package preserves the same developer experience — users import one package and access all public APIs. The codebase is ~2,200 LOC, small enough that a single package remains clean and navigable.

---

## 2. `ScanToken.Depth` — `int` Instead of `float64`

**Decision**: Use `int` with `math.MaxInt32` as the sentinel for "infinite" depth.

**Rationale**: The JS source uses `Infinity` (a float) for globstar tokens. In Go, converting `math.Inf(1)` to `int` is undefined behavior. We use `math.MaxInt32` as a safe, large sentinel that behaves correctly in integer arithmetic.

---

## 3. RE2 Regex Engine Compatibility

**Decision**: Replace JS-specific regex patterns that use backreferences with RE2-compatible alternatives.

**Rationale**: Go's `regexp` package uses the RE2 engine, which does not support backreferences (`\1`, `\3`, etc.) or lookahead/lookbehind assertions. The `REGEX_SPECIAL_CHARS_BACKREF` pattern was simplified to match the core behavior without backreferences.

---

## 4. Zero External Dependencies

**Decision**: No third-party Go modules.

**Rationale**: The original JS `picomatch` has zero dependencies — it's a core selling point. We preserve this property in the Go port, using only the Go standard library (`regexp`, `strings`, `math`, `unicode`).

---

## 5. Options Struct Instead of Variadic Maps

**Decision**: Use a strongly-typed `Options` struct with named fields.

**Rationale**: JS `picomatch` accepts an options object `{dot: true, nocase: true, ...}`. In Go, a struct with typed fields provides compile-time safety, IDE autocompletion, and clear documentation — a strict improvement over the JS pattern. Passing `nil` for defaults is idiomatic Go.

---

## 6. `RemoveBackslashes` — Regex-Based vs Manual Loop

**Decision**: Use `REGEX_REMOVE_BACKSLASH` (compiled regex) for backslash removal.

**Rationale**: The JS source uses a regex replace. While a manual byte loop would be slightly faster in Go, using the same regex approach ensures behavioral parity with the original and reduces the risk of edge-case divergence. Performance can be optimized later if benchmarks warrant it.

---

## 7. Windows Path Support via `GlobChars`

**Decision**: Maintain separate `PosixChars` and `WindowsChars` structs with platform-specific regex tokens.

**Rationale**: Mirrors the original JS `constants.js` design exactly. The `GetGlobChars(win32 bool)` function selects the appropriate token set at runtime, enabling cross-platform glob matching without `runtime.GOOS` checks in hot paths.

---

## 8. Incremental PR Workflow

**Decision**: Build the port incrementally via small, focused pull requests (options → constants → utils → scanner → parser → compiler → matcher).

**Rationale**: The hackathon rules require "real, incremental commit history." Small PRs also enable parallel code review by team members and provide clear git bisect targets if regressions occur.
