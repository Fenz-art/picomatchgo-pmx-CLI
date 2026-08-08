# PMX Engineering Rules

PMX is a Picomatch-compatible Go runtime, a command-line validation product, and a Foundry dashboard. Preserve the boundary between the runtime, CLI, dashboard API, and browser UI. Custom agent and skill definitions live under `.codex/` and are the source of truth for automated validation workflows.

## Required workflow

1. Inspect the affected runtime, command, route, fixture, and documentation before editing.
2. Prefer structured PMX output (`pmx doctor --json`, `pmx agent inspect --json`, and `pmx agent check --json`) over parsing human terminal text.
3. Run the narrowest relevant test after a change, then the repository regression surface before handoff.
4. Treat a feature as verified only after executing it and observing its output. A static success label is never evidence.
5. Keep Foundry status derived from API command exit code, stdout/stderr, duration, and structured payloads. Do not add simulated workflows, fabricated logs, or hard-coded PASS states.

## Boundaries and quality gates

- Keep scan, parse, compile, and match independently testable.
- CLI commands must provide useful errors and non-zero exits for invalid input or failed execution.
- Dashboard command execution must use argument arrays (`execFile`), an allowlist, a repository-root working directory, timeouts, and JSON errors.
- GitHub tokens remain server-only. Never expose them through client components or API responses.
- Tests must cover both successful and failure paths where realistic. Do not weaken a test merely to hide a regression.
- Format Go with `gofmt`; keep `go vet ./...`, `go test ./...`, and `go test -race ./...` green. Keep dashboard lint/build green when Node tooling is available.

## Submission evidence

Maintain `docs/architecture.md`, this file, and `AGENTS_AND_SKILLS.md`. The repository must contain a usable custom agent and custom skill under `.codex/`, and CI jobs must execute real checks.
