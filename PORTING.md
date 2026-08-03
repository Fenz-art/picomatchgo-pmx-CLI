# Porting Notes

## Goal

The goal of this repository is to provide a faithful Go port of the core behavior of picomatch while staying idiomatic for Go and avoiding external dependencies.

## Porting strategy

- Preserve the public behavior of the original library where possible.
- Keep the implementation organized into scanner, parser, matcher, and utility layers.
- Favor pure Go standard-library primitives over JavaScript-specific assumptions.
- Cover advanced features such as extglobs, brace expansion, dotfiles, and POSIX classes.

## Compatibility notes

- Path separators are normalized by the scanner and parser helpers to match the original library across POSIX and Windows-style input.
- Dotfile semantics are handled in the matcher layer to preserve upstream behavior for patterns that do not explicitly start with a dot.
- The implementation uses Go regexp syntax and RE2-compatible constructs, which shapes some edge-case behavior around lookarounds and complex backtracking.
- The parser keeps the public API simple by exposing Scan, Parse, MakeRe, CompileRe, and IsMatch while preserving the same feature set for extglobs, braces, globstars, and POSIX classes.
