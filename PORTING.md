# Porting Notes

## Goal

The goal of this repository is to provide a faithful Go port of the core behavior of picomatch while staying idiomatic for Go and avoiding external dependencies.

## Porting strategy

- Preserve the public behavior of the original library where possible.
- Keep the implementation organized into scanner, parser, matcher, and utility layers.
- Favor pure Go standard-library primitives over JavaScript-specific assumptions.
- Cover advanced features such as extglobs, brace expansion, dotfiles, and POSIX classes.

## Compatibility notes

- Path separators are normalized by the scanner and parser helpers.
- Dotfile semantics are handled in the matcher layer to preserve upstream behavior.
- The implementation uses Go regexp syntax and RE2-compatible constructs.
