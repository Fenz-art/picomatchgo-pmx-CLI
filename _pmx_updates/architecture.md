# PMX Architecture

PMX is a dependency-free Go implementation of core Picomatch glob behavior, wrapped by an executable developer-validation CLI and a Next.js Foundry dashboard. The system is intentionally layered so each result can be traced to real execution.

```text
Pattern/input
  -> Go runtime: Scan -> Parse -> CompileRe -> IsMatch
  -> PMX CLI: match/scan/parse/explain/validate/compat/doctor/agent/ci
  -> JSON contracts: doctor and ADLC reports
  -> Foundry validation API
  -> Foundry UI

GitHub Actions workflow -> GitHub API routes -> Foundry workflow UI
```

## Runtime

The root Go module supplies the public runtime. `scan_impl.go` produces structural pattern metadata; `parse_impl.go` converts supported glob syntax to RE2-compatible regex source; `matcher_impl.go` compiles and evaluates it. `options.go`, `types.go`, `constants.go`, and `utils_impl.go` hold the shared public configuration, data contracts, constants, and helpers. This separation is deliberate: scan, parse, compile, and match remain independently testable.

`cmd/wasm` builds the same runtime for the browser. The Foundry Engine Lab loads `public/picomatch.wasm`; failures to load or execute are surfaced in the UI rather than treated as passing execution.

## PMX CLI and diagnostics

`cmd/pmx` is the command boundary. Its public commands execute the runtime or named repository checks. `doctor` detects project configuration and emits human or JSON reports. The diagnostic contract includes a code/id, severity, message, relevant file where known, and suggested action.

The ADLC commands are:

- `pmx agent inspect --json`: project snapshot plus diagnostics.
- `pmx agent check --json`: canonical `{version, result, diagnostics, checks, next_actions}` gate.

Agent checks spawn the current executable and run the same doctor, validation, compatibility, CI, and regression paths that a developer invokes. A failed child command becomes a failed check; the result is not inferred from documentation or UI state.

`pmx ci --json` executes local format, vet, unit, race, CLI, compatibility, doctor, and regression checks. It is a local validation report, not a GitHub Actions client. Live workflow status and logs belong to Foundry's GitHub integration.

## Foundry

The Next.js application in `dashboard/` has two execution paths:

1. `/api/validation` accepts only allowlisted validation IDs and calls `go run ./cmd/pmx` with argument arrays, a fixed repository root, timeout, stdout/stderr capture, exit code, duration, and optional parsed JSON. The UI renders those returned results, including pass, warn, fail, and errors.
2. `/api/workflows` and nested run/job/artifact/log routes use a server-side GitHub token to dispatch and query GitHub Actions. The browser never receives the token. Missing credentials and GitHub API errors are returned as JSON errors and displayed as failures/unavailable state.

The dashboard must not represent local CLI results as a GitHub run, or represent unavailable GitHub data as a pass.

`cmd/dashboard` is retained only as a legacy static UI fixture for its package tests. It carries an explicit demo marker and must not be used as submission evidence or deployed as Foundry.

## CI/CD and deployment

`.github/workflows/ci.yml` runs formatting, linting, vet, unit, race, CLI smoke commands, doctor, compatibility, regression, fuzzing, benchmarks, and dashboard lint/build. It is triggered on `main` push, pull requests to `main`, and `workflow_dispatch`. The doctor report is uploaded as an artifact.

The dashboard is deployable as a standard Next.js service. Configure `FOUNDRY_GITHUB_REPOSITORY` as `owner/repo` and provide `FOUNDRY_GITHUB_TOKEN` only to the server environment. The Go toolchain must be available to deployments that use `/api/validation`.

## Repository evidence

Architecture and engineering rules are in `docs/architecture.md` and `AGENTS.md`. The custom PMX Validation Agent and PMX Doctor/ADLC skill are documented in `AGENTS_AND_SKILLS.md` and stored under `.codex/` in this repository.
