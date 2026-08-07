# Picomatch Go → Track B Execution Plan
### `pmx` — the glob doctor, built on the runtime you already shipped

---

## 1. EXECUTIVE VERDICT

**"Picomatch Doctor" / `pmx` is the right idea — but only the diagnostic half of it, and only if you delete the fake CI simulator.**

Two hard findings from reading the actual repo:

- `ci.yml` is one job (`quality`) triggered by `push`/`pull_request`. There is **no `workflow_dispatch`**, so nothing in a web dashboard can currently trigger it on click. This must be added — it's small, but it's not optional, it's the whole point of "Run CI."
- The README documents a "Live Workflows" tab that is explicitly a **simulator**: "Interactive CI/CD pipeline simulator... Trigger Integration Pipeline... 8 validation steps execute sequentially with real-time terminal log output." That is fabricated execution. Both hackathon docs you gave me explicitly forbid this ("fake workflow execution," "fake CI," "fake logs"). If a judge inspects network traffic and sees no request leave the browser, you don't lose points — you lose the entry checkpoint on green, real CI, because the dashboard actively demonstrates you're comfortable faking it.

Everything else in your existing stack — Scan, Parse, MakeRe, IsMatch, the fuzz targets, the bench suite, the WASM build — is real, reusable, and exactly the right foundation. You don't need a new idea. You need to point the existing dashboard at existing GitHub Actions instead of a `setTimeout` loop, and put a thin, honest CLI on top of the matcher internals you already wrote.

---

## 2. FINAL PRODUCT

**`pmx` — Picomatch Doctor**, backed by the **Picomatch Go Foundry**, a real (not simulated) validation control plane for the same repo.

**ONE-LINE PITCH:** A CLI and web control plane that tells a developer *exactly why* a glob pattern did or didn't match — using the real scanner/parser/matcher internals — and proves the engine is trustworthy by running real GitHub Actions CI, fuzzing, and JS-compatibility checks on demand.

**DEVELOPER PROBLEM:** "Why didn't `src/**/*.ts` match `src/index.tsx`?" Developers debug globs by trial-and-error edits and reruns of whatever tool consumes the pattern (bundler ignore rules, CI path filters, `.gitignore`-style configs, test globs). There is no standard tool that shows the *stages* — scan → parse → compiled regex → per-segment match/no-match — the way `EXPLAIN`-style tools exist for SQL or regex debuggers exist for `regex101`. Picomatch Go's pipeline is naturally that debugger; it just isn't exposed as one yet.

**WHY PICOMATCH GO IS THE FOUNDATION:** `Scan`, `Parse`, and `MakeRe` already expose exactly the intermediate state a diagnosis needs — base path, glob flag, brace/extglob flags, AST, compiled regex source. `pmx explain`/`pmx diagnose` are almost entirely formatting layers over functions that already exist in `scan_impl.go` and `parse_impl.go`. This is a case where the hard engineering (the port itself, plus fuzzing/compat proving it's correct) is done; the productization is thin.

**WHY THIS FITS TRACK B:** It's a CLI + GitHub Action + web app pointed at a real developer workflow (debugging why a glob rule in CI, a bundler config, or a test filter didn't do what was expected), matching the track's own definition almost word for word.

**WHY A JUDGE WILL CARE:** Judges see a lot of "we ported X to Go" submissions where the demo is `go test ./...` passing. This one lets the judge type a pattern that's *wrong*, watch the tool point at the exact AST node or scan flag responsible, then click a button that kicks off the actual GitHub Actions run validating that answer — with a real run ID, real logs, real pass/fail. That's the difference between "trust us, it's correct" and "here's the run, go look."

---

## 3. WHY THIS WINS TRACK B

Judges are told to look for: fits a real workflow, integrates as CLI/Action/web app, reliable enough to actually use. `pmx` clears all three without invention:

- **CLI**: `pmx` ships as a normal Go binary (`cmd/pmx`), works offline, no network dependency for the diagnostic commands.
- **GitHub Action**: the same binary runs as a repo-local composite action (`pmx validate` in CI, or `pmx diagnose` on a PR that changes glob-driven config) — trivial to add since it's already a Go module with a `cmd/` convention.
- **Web app**: the Foundry becomes the observability layer over the *same* validate/fuzz/bench/compat logic the CLI runs locally, wired to real GitHub Actions runs instead of a client-side animation.

Nothing here requires a new domain concept, a database, or an LLM. The "productivity tool" is a debugger for a thing your team already built a correct implementation of — which is the strongest possible story for "we understood the engineering problem," because the debugger's correctness *is* the port's correctness.

---

## 4. CORE USER PROBLEM

A developer has a glob-driven config — CI path filters, bundler include/exclude, test file discovery, `.gitignore`-style ignore rules — and a pattern that isn't matching what they expect. Today they:

1. Guess a pattern.
2. Run whatever consumes it (webpack, jest, a CI workflow).
3. Get a boolean/silent failure with no explanation.
4. Repeat.

`pmx` collapses that loop: give it the pattern and the path, get back *which stage* rejected it and *why*, deterministically, in under 50ms, with zero network calls.

---

## 5. FIXED EXECUTION FLOW

This is the one flow the whole product — CLI, Action, and Foundry — is built around. Do not deviate from it in the demo.

```
Developer has a glob that "should" match a path but doesn't
        ↓
$ pmx diagnose "<pattern>" "<path>"
        ↓
pmx calls Scan()   → shows base/glob/brace/extglob flags
pmx calls Parse()  → shows AST + which token broke assumptions
pmx calls MakeRe() → shows compiled regex
pmx calls IsMatch()→ shows true/false + the exact failing regex anchor
        ↓
pmx prints a one-line verdict: WHY it didn't match
        ↓
Developer fixes the pattern, re-runs `pmx diagnose`, gets a match
        ↓
Developer runs `pmx validate` (or clicks "Run CI" in the Foundry)
        ↓
Real GitHub Actions workflow_dispatch fires on ci.yml
        ↓
Foundry polls the real Actions API: queued → in_progress → per-job status
        ↓
Judge watches real logs stream from a real run ID
        ↓
Foundry shows pass/fail pulled from the real conclusion field
```

Everything left of "Developer runs `pmx validate`" is 100% local, deterministic, offline. Everything right of it is 100% real GitHub Actions, no simulation layer.

---

## 6. P0 MVP

| Feature | Dev value | Judge value | Effort | Already exists? |
|---|---|---|---|---|
| `pmx match` | quick match test | baseline | trivial | Yes — wraps `IsMatch` |
| `pmx explain` | shows Scan+Parse output | shows internals are real, not a black box | small | Yes — wraps `Scan`/`Parse`, needs formatting only |
| `pmx diagnose` | root-causes a non-match | the "wow" command | medium | Mostly — needs a diff between Parse AST and match trace, new logic ~150-250 LoC |
| `pmx validate` (local) | runs fmt/vet/test/race locally, same commands as CI | proves CLI and CI use identical logic | small | Yes — literally shells the same commands `ci.yml` runs |
| `workflow_dispatch` on `ci.yml` | — | required for any real "Run CI" button | small | No — must add trigger + `github.token`-safe inputs |
| Foundry "Run CI" wired to real Actions API | — | required to not be fake | medium | No — replace the simulator's `setTimeout` sequence with real `POST /actions/workflows/ci.yml/dispatches` + polling `GET /actions/runs` |
| Foundry live log tail | — | the moment judges remember | medium | No — pull `GET /actions/jobs/{id}/logs`, stream to UI |

**P1 — only if time remains**

| Feature | Dev value | Judge value | Effort | Already exists? |
|---|---|---|---|---|
| `pmx compat` | proves JS parity on demand | strong "not just a port" signal | medium | Partial — README claims a differential fuzzer exists; verify it's real, wire a small fixed corpus (30–50 patterns) run through both engines in a dedicated CI job |
| `pmx bench` | shows perf numbers | credibility | small | Yes — `bench_test.go` already exists, just needs `-json` parsing + a stored baseline file to diff against |
| Foundry fuzz evidence panel | shows exec count, corpus growth, crashers found | proves fuzzing isn't decorative | small-medium | Partial — `FuzzScan`/`FuzzParse`/`FuzzIsMatch` exist, need `go test -fuzz` output parsed into structured JSON |
| GitHub composite Action wrapping `pmx validate` | lets other repos use it | "genuinely usable" signal | small | No, but trivial given `cmd/pmx` exists |

**P2 — do not build**

- AI/LLM explanation layer ("ask AI why this didn't match") — the deterministic explanation from Scan/Parse *is* the differentiator; an LLM guessing on top of a correct symbolic engine is a downgrade, not an upgrade.
- User accounts, saved patterns, history, multi-tenant anything.
- A generic "CI replacement" — you already have GitHub Actions; don't reimplement a scheduler.
- Coverage percentages beyond what `go test -coverprofile` actually reports.
- Any second language runtime in the demo path (no live Node process needed at runtime — precompute compat fixtures in CI instead, see §11).

---

## 7. CLI

Six commands, all in `cmd/pmx`, all thin wrappers over existing exported functions plus one new diagnostic pass.

### `$ pmx match <pattern> <path> [flags]`
```
$ pmx match "src/**/*.ts" "src/index.tsx"
NO MATCH
```
Invokes `picomatch.IsMatch`. Flags mirror `Options` (`--dot`, `--nocase`, `--posix`, `--windows`). Why a developer needs it: fastest possible sanity check, scriptable in a shell pipeline (exit code 0/1). Why a judge should see it: establishes the baseline before `diagnose` shows its work.

### `$ pmx explain <pattern>`
```
$ pmx explain "src/**/*.{ts,tsx}"
SCAN
  base       = src
  glob       = **/*.{ts,tsx}
  isGlob     = true
  isBrace    = true
  isGlobstar = true

PARSE (AST, top-level nodes)
  segment: **        → globstar (matches zero or more path segments)
  segment: *          → wildcard (no dot-files unless --dot)
  segment: {ts,tsx}    → brace-set, expands to 2 alternatives

COMPILED REGEX
  ^(?:src(?:/|$)(?:.*\/)?[^/]*\.(?:ts|tsx))$
```
Invokes `Scan` then `Parse` then `MakeRe`, printing the real intermediate state — no separate "explain engine." Developer need: understand a pattern before trusting it in config. Judge need: proves the internals aren't a black box; this is the same data structures the matcher actually runs on, not a re-derived summary.

### `$ pmx diagnose <pattern> <path> [flags]`
```
$ pmx diagnose "src/**/*.ts" "src/index.tsx"
VERDICT: NO MATCH

  pattern compiled to: ^(?:src/(?:.*\/)?[^/]*\.ts)$
  path                : src/index.tsx

  Extension check failed:
    compiled group requires suffix ".ts" (literal, non-optional)
    "index.tsx" ends in ".tsx", not ".ts"

  Nearest fix:
    pattern "src/**/*.{ts,tsx}" WOULD match this path
    (verified: re-running IsMatch with suggested pattern → MATCH)
```
This is the one new component: after a non-match, `pmx` walks the compiled regex's literal segments against the input to identify which literal/anchor first diverges, and — deterministically, not via an LLM — tries a small fixed set of mechanical relaxations (turn a literal extension into a brace set already implied by sibling files in the same dir if `--path-context` is given; loosen an extglob boundary) and *actually re-runs `IsMatch`* to confirm the suggestion before printing it. Never print an unverified suggestion. Developer need: this is the whole point of the tool. Judge need: this is the moment that makes it "not just a wrapper" — the fix suggestion is proven live, not guessed.

### `$ pmx validate [--fast]`
```
$ pmx validate
✓ gofmt            0.4s
✓ go vet            1.1s
✓ go test ./...     3.8s
✓ go test -race     9.2s
  (skipped: fuzz, bench — use --full or Foundry "Run CI")
ALL CHECKS PASSED
```
Shells the exact same commands defined in `ci.yml`'s `quality` job. `--full` also runs the fuzz/bench steps locally. Developer need: pre-push confidence, identical to what CI will do. Judge need: literal proof the CLI and the CI config aren't two different stories — same commands, shown side by side if you open `ci.yml` in the demo.

### `$ pmx compat <pattern>` (P1)
```
$ pmx compat "!(*.test).js"
Go     : matches ["a.js"], rejects ["a.test.js"]  (12 fixture paths)
JS ref : matches ["a.js"], rejects ["a.test.js"]  (12 fixture paths)
PARITY: 12/12
```
Runs the pattern against a fixed fixture corpus checked into the repo, whose *expected* results were generated once (offline, in CI) by running the real `micromatch/picomatch` JS package against the same fixtures. Developer need: confidence when migrating from JS picomatch. Judge need: concrete, non-inflated compatibility number tied to a visible fixture file, not a claimed percentage.

### `$ pmx bench [--baseline]` (P1)
```
$ pmx bench
IsMatch   412 ns/op   128 B/op   4 allocs/op   (baseline: 415 ns/op, -0.7%)
Scan      244 ns/op    64 B/op   2 allocs/op   (baseline: 243 ns/op, +0.4%)
Parse     591 ns/op   256 B/op   8 allocs/op   (baseline: 588 ns/op, +0.5%)
```
Parses `go test -bench=. -benchmem -json` output and diffs against a baseline JSON committed to the repo (regenerated on tagged releases only). Developer need: catch perf regressions before merge. Judge need: real numbers, real diff, not narrative claims.

---

## 8. FOUNDRY

Cut the tab count. Five surfaces, not ten:

**Navigation:** `Diagnose` · `Validate` · `Compatibility` · `Fuzzing` · `Benchmarks`

Drop "Artifacts," "Release Readiness," and "Validation Matrix" as standalone tabs — fold artifact links into the Validate run detail view, and fold "Validation Matrix" into the Diagnose page's explain panel. Fewer tabs a judge has to click through in a 4-minute demo.

**Primary buttons, exact labels:**
- `Run CI` — dispatches `ci.yml` with `workflow_dispatch`
- `Run Fuzz` — dispatches with `input: mode=fuzz`
- `Run Benchmarks` — dispatches with `input: mode=bench`
- `Run Compatibility` — dispatches with `input: mode=compat`
- `View Run on GitHub` — deep link to the real Actions run URL (always present, always real — this single link is your strongest anti-fake proof)

**States, exact copy:**
- Empty: `No run yet. Click "Run CI" to validate the current commit on GitHub Actions.`
- Queued: `QUEUED — waiting for a GitHub-hosted runner`
- Running: `RUNNING — <job name> (<elapsed>)`
- Passed: `PASSED — run #<id> · <duration> · <sha short>`
- Failed: `FAILED at <job/step name> — <first error line from real log>`

**Diagnose tab:** the pattern/path input from `pmx diagnose`, running WASM-compiled `Scan`/`Parse`/`MakeRe`/`IsMatch` in-browser (you already build `cmd/wasm/main.go` for this — reuse it, don't rebuild a server endpoint). This is real: it's the literal Go engine compiled to WASM, executing client-side, deterministic, no server round trip. That's a legitimate "real execution" claim distinct from the CI panel.

**Validate/Fuzz/Bench/Compat tabs:** identical shape — a run history list (real run IDs from the Actions API), a detail view with per-job status pulled from `GET /repos/{owner}/{repo}/actions/runs/{run_id}/jobs`, and a log pane streaming from `GET .../jobs/{job_id}/logs` (this endpoint returns a redirect to a downloadable zip/text; fetch server-side and stream lines to the client — do not attempt to fake incremental streaming, just poll and append).

---

## 9. REAL CI/CD EXECUTION

```
Foundry "Run CI" button
   ↓
Next.js API route (server-side, holds the token)
   ↓
POST /repos/{owner}/{repo}/actions/workflows/ci.yml/dispatches
   body: { ref: "main", inputs: { mode: "quality" } }
   ↓
GitHub Actions schedules a run
   ↓
Foundry polls GET /repos/{owner}/{repo}/actions/runs?event=workflow_dispatch
   (match by created_at proximity + head_sha, since dispatch has no run_id in response)
   ↓
Once run_id resolved: poll GET /actions/runs/{run_id} → status/conclusion
                       poll GET /actions/runs/{run_id}/jobs → per-job status/steps
   ↓
On completion: fetch logs, fetch artifacts list via GET /actions/runs/{run_id}/artifacts
   ↓
Foundry renders final state
```

**workflow_dispatch inputs** (add to `ci.yml`):
```yaml
on:
  push: { branches: [main] }
  pull_request: { branches: [main] }
  workflow_dispatch:
    inputs:
      mode:
        description: "quality | bench | fuzz | compat"
        default: "quality"
        type: choice
        options: [quality, bench, fuzz, compat]
```
Then split the current single `quality` job into `quality`, `bench`, `fuzz`, `compat` jobs, each gated with `if: github.event.inputs.mode == 'X' || github.event_name != 'workflow_dispatch'` (so push/PR still runs everything, but a dashboard-triggered run only runs the requested job — this is what makes "Run Benchmarks" show a fast, focused run instead of the whole suite).

**Authentication:** a fine-grained GitHub PAT (or a GitHub App installation token, if you want it to look sharper) with `actions:write` + `actions:read` scoped to only this repo, stored as a Vercel encrypted env var, used exclusively server-side in the API route. Never exposed to the client. State this explicitly in the demo — "the token never touches the browser" is a real, checkable security boundary, and mentioning it preempts a judge question.

**Run IDs / polling:** dispatch responses don't return a run ID (a known GitHub API limitation) — poll the runs list filtered by `head_sha` immediately after dispatch, typically resolves within 1–3s.

**Failure handling:** on `conclusion: "failure"`, fetch the failing job's steps, find the first `conclusion: "failure"` step, pull its log slice, surface that as the headline error rather than a generic "failed."

**Cancellation/re-run:** `POST /actions/runs/{run_id}/cancel` and `POST /actions/runs/{run_id}/rerun` — both are real endpoints, wire both if time allows (cheap, high demo value: "if this fails, watch us re-run it live").

**Security boundary:** dashboard is read/dispatch only — no code from dashboard input is ever committed, executed, or interpolated into workflow YAML. The only user-controlled input reaching GitHub is the `mode` enum, which is a `choice` type input GitHub itself validates against the fixed option list — no injection surface.

---

## 10. VALIDATION ENGINE

| NAME | PURPOSE | COMMAND | INPUT | EXPECTED | OUTPUT | FAILURE MEANS |
|---|---|---|---|---|---|---|
| Format | canonical gofmt | `gofmt -l .` | repo `.go` files | empty output | file list or empty | someone committed unformatted Go |
| Vet | suspicious constructs | `go vet ./...` | repo | exit 0 | diagnostics | likely real bug (nil deref, bad format verb, etc.) |
| Unit | correctness | `go test ./...` | `*_test.go` | all pass | pass/fail per package | scanner/parser/matcher regression |
| Race | concurrency safety | `go test -race ./...` | same, `-race` | no data races | pass/fail | matcher/parser not safe for concurrent reuse (relevant since `MakeRe` compiles regexes that may be cached) |
| Fuzz | crash/mismatch discovery | `go test -fuzz=FuzzScan -fuzztime=15s .` (×3 targets) | seed corpus + generated inputs | no crashers | executions, new corpus entries, crashers | scanner/parser panics or hangs on adversarial input |
| Bench | perf regression | `go test -bench=. -benchmem -run=^$ ./...` | benchmark inputs in `bench_test.go` | within tolerance of baseline | ns/op, B/op, allocs/op | perf regression on hot path |
| Compat | JS parity | custom `compat_test.go` against fixture JSON | fixed pattern/path corpus | Go result == JS-precomputed result | parity count | behavioral divergence from `micromatch/picomatch` |

Keep this list exactly this size. Do not add lint-rule-by-lint-rule breakdowns or per-function coverage tables — that's checklist theater, not engineering story.

---

## 11. COMPATIBILITY

Smallest credible design, no live Node process at demo time:

1. Build a fixed fixture file: 30–50 `(pattern, path, options)` triples chosen to exercise globstar, braces, extglobs, negation, dotfiles, POSIX classes, Windows paths — the exact feature list your README already claims to support.
2. **Once**, offline (or as a one-time CI job, not in the demo hot path), run those fixtures through the real `micromatch/picomatch` npm package in Node and record `expected.json`.
3. Commit `expected.json` to the repo.
4. `compat_test.go` (Go) loads the fixture file, runs each triple through `picomatch.IsMatch`, and diffs against `expected.json`. This is what `pmx compat` and the Foundry's Compatibility tab both call.
5. Report the real fraction: `47/50` if that's the truth. If it's not 100%, that's fine — say so, list the 3 divergent patterns, and note them as known JS-parity gaps. A judge trusts "47/50, here are the 3 gaps and why" far more than an unverifiable "100% compatible" claim, and it directly defuses "how do you prove compatibility?"

Do not build a live Node subprocess call at runtime in the CLI or the browser — that reintroduces a Node dependency into a "zero-dependency Go" pitch and is a demo-day fragility risk (network/npm install failures live on stage). Precomputed fixtures avoid that entirely while still being real, checkable, and regenerable.

---

## 12. FUZZ + BENCHMARKS

**Fuzzing — show real evidence, not a pass/fail pill:**
```
FuzzScan       15s   1,204,332 execs   0 crashers   3 new corpus entries
FuzzParse      15s     812,109 execs   0 crashers   1 new corpus entry
FuzzIsMatch    15s     640,221 execs   0 crashers   0 new corpus entries
```
`go test -fuzz=X -fuzztime=15s` prints exec counts and corpus growth to stdout/stderr already; the CI job just needs to capture that and the API route parses it into these fields. If a crasher is found, `testdata/fuzz/FuzzX/<hash>` gets a new file — link directly to it in the run artifacts, that's the single most convincing "this is real" artifact you can show a judge, since it's a literal failing input Go found on its own. Do not claim coverage percentages — Go's native fuzzer doesn't expose coverage % without extra tooling (`go tool cover` doesn't apply here the same way); state plainly "coverage % not available for native fuzz targets" rather than inventing a number.

**Benchmarks — real numbers only, from `bench_test.go`, never the sample numbers currently pasted in the README** (those "Apple M2, Go 1.21" numbers in the README are static/stale — replace with numbers generated by the actual CI run, or clearly label them "illustrative, regenerate via `pmx bench`"). Comparison baseline: commit a `baseline.json` at each tagged release (§17); `pmx bench` and the Foundry both diff the current run against the last tagged baseline, not against the previous run, so numbers are stable and don't flap on noisy runner hardware.

---

## 13. ARCHITECTURE

```
Picomatch Go core (types.go, options.go, scan_impl.go, parse_impl.go, matcher_impl.go)
      ↓                                   ↓
   cmd/pmx (CLI)                    cmd/wasm (WASM build)
      ↓                                   ↓
.github/workflows/ci.yml          dashboard/ (in-browser Diagnose tab,
(quality/bench/fuzz/compat jobs)   calls WASM directly, no server hop)
      ↓
GitHub Actions REST API
      ↓
dashboard/app/api/* (Next.js API routes — server-side token, dispatch + poll)
      ↓
Foundry UI (Validate/Fuzz/Bench/Compat tabs — real run data only)
```

Two independent "real execution" paths by design: WASM in-browser for instant local diagnosis (no server, no token, no GitHub dependency — works even if Actions is down), and GitHub Actions for anything that needs CI-grade reproducibility (race detector, fuzz, full test matrix). Never let the CI path fall back to a client-side simulation "in case the API is slow" — if the dispatch or poll fails, show a real error state, not a fake progress bar.

---

## 14. REPOSITORY STRUCTURE

Adapted from what exists — no destructive refactor, additive only:

```
port-mortem-picomatch-go/
├── .github/workflows/ci.yml        # MODIFIED: add workflow_dispatch, split into 4 jobs
├── cmd/
│   ├── wasm/main.go                 # existing
│   └── pmx/                         # NEW: CLI entrypoint
│       ├── main.go
│       ├── cmd_match.go
│       ├── cmd_explain.go
│       ├── cmd_diagnose.go
│       ├── cmd_validate.go
│       ├── cmd_compat.go
│       └── cmd_bench.go
├── compat/                          # NEW: compatibility fixtures + generator
│   ├── fixtures.json
│   ├── expected.json                # precomputed from real JS picomatch, committed
│   └── generate/                    # one-off Node script, NOT part of build/runtime path
│       └── generate.js
├── compat_test.go                   # NEW: Go side of §11
├── baseline.json                    # NEW: bench baseline, regenerated on tagged releases
├── scan_impl.go / parse_impl.go / matcher_impl.go / types.go / options.go   # unchanged
├── bench_test.go / fuzz_test.go / picomatch_test.go / scan_test.go          # unchanged
├── dashboard/
│   ├── app/
│   │   ├── api/
│   │   │   ├── ci/dispatch/route.ts     # NEW: POST workflow_dispatch
│   │   │   ├── ci/status/route.ts       # NEW: GET run + jobs status
│   │   │   └── ci/logs/route.ts         # NEW: GET job logs
│   │   ├── diagnose/page.tsx            # RENAMED from playground, WASM-backed
│   │   ├── validate/page.tsx            # REPLACES "Live Workflows" simulator
│   │   ├── compat/page.tsx
│   │   ├── fuzz/page.tsx
│   │   └── bench/page.tsx
│   └── public/picomatch.wasm / wasm_exec.js  # unchanged
├── AGENTS.md                        # NEW: entry checkpoint #2
├── AGENTS_AND_SKILLS.md             # NEW: entry checkpoint #4
├── docs/architecture.md             # existing, update to reflect pmx
├── docs/prd.md                      # NEW: user stories + acceptance criteria
└── test/e2e/                        # NEW: Playwright specs
```

---

## 15. IMPLEMENTATION PLAN

**STEP 1 — Add `workflow_dispatch` and split `ci.yml` into 4 jobs**
Files: `.github/workflows/ci.yml`
Implementation: add `workflow_dispatch` trigger with `mode` choice input; split the current single `quality` job's steps into `quality` (fmt/lint/vet/test/race), `bench`, `fuzz`, `compat` jobs, each gated by the `mode` input on manual dispatch, all running unconditionally on push/PR.
Test: push a commit, confirm all 4 jobs still run and pass exactly as the current single job did; manually dispatch with `mode=bench`, confirm only that job runs.
Definition of done: a green Actions run exists from a manual `workflow_dispatch`, visible in the Actions tab.

**STEP 2 — Build `cmd/pmx` with `match`, `explain`, `validate`**
Files: `cmd/pmx/*.go`
Implementation: thin CLI (stdlib `flag` or minimal arg parsing, no dependency needed) calling `picomatch.Scan/Parse/MakeRe/IsMatch` and shelling the exact commands from the `quality` job for `validate`.
Test: unit tests per subcommand with fixed pattern/path pairs and golden output.
Definition of done: `go build ./cmd/pmx && ./pmx match "*.js" "a.js"` works from a clean checkout.

**STEP 3 — Build `pmx diagnose`**
Files: `cmd/pmx/cmd_diagnose.go`, new internal package `internal/diagnose`
Implementation: on non-match, walk `ParseState`'s literal/anchor segments against the input, identify first divergence, attempt a small fixed set of mechanical pattern relaxations, re-verify each candidate with a real `IsMatch` call before printing.
Test: table-driven tests covering globstar-vs-single-star, missing brace alternative, dotfile exclusion, case mismatch.
Definition of done: `pmx diagnose "src/**/*.ts" "src/index.tsx"` prints a correct, re-verified suggestion.

**STEP 4 — Compatibility fixtures**
Files: `compat/generate/generate.js`, `compat/fixtures.json`, `compat/expected.json`, `compat_test.go`, `cmd/pmx/cmd_compat.go`
Implementation: write ~40 fixture triples covering the README's claimed feature list; run once via real `micromatch/picomatch` in Node to produce `expected.json`; write `compat_test.go` to diff Go output against it; add `pmx compat`.
Test: `go test -run TestCompat ./...` passes (or documents real divergences).
Definition of done: fixture file + expected file both committed, test green (or divergences documented in `docs/compat-gaps.md`).

**STEP 5 — Bench baseline + `pmx bench`**
Files: `baseline.json`, `cmd/pmx/cmd_bench.go`
Implementation: run `go test -bench=. -benchmem -json`, parse, compare to `baseline.json`, print deltas.
Test: run twice locally, confirm delta math is correct and near-zero on identical hardware.
Definition of done: `pmx bench` prints real ns/op numbers with a delta line, no hardcoded numbers anywhere.

**STEP 6 — Next.js API routes for real CI dispatch/poll**
Files: `dashboard/app/api/ci/dispatch/route.ts`, `.../status/route.ts`, `.../logs/route.ts`
Implementation: server-side fetch to GitHub REST API using a Vercel env var token; dispatch, resolve run_id via head_sha match, poll status/jobs, fetch logs on completion.
Test: manual — click dispatch locally, confirm a real run appears in the Actions tab within seconds.
Definition of done: dashboard can dispatch and display a real run end-to-end with no client-visible token.

**STEP 7 — Replace the simulator with real UI**
Files: `dashboard/app/validate/page.tsx` (was the "Live Workflows" simulator), `compat/`, `fuzz/`, `bench/` pages
Implementation: delete the `setTimeout`-driven step sequence entirely; wire buttons to Step 6's API routes; render real state machine (queued/running/passed/failed) from polled data.
Test: full click-through against a live repo run.
Definition of done: no `setTimeout`, no hardcoded step list, no client-side "log" strings anywhere in this code path — grep the diff to confirm.

**STEP 8 — WASM-backed Diagnose tab**
Files: `dashboard/app/diagnose/page.tsx`
Implementation: reuse existing `picomatch.wasm`/`wasm_exec.js`, call `picomatchScan`/`picomatchParse`/`picomatchIsMatch`/`picomatchCompile` directly from the pattern/path inputs, render the same structure as `pmx explain`/`diagnose`.
Test: manual, cross-check against CLI output for the same inputs (they should be byte-identical since it's the same compiled code).
Definition of done: in-browser diagnosis matches CLI diagnosis exactly, no server round-trip.

**STEP 9 — Hackathon entry checkpoints**
Files: `AGENTS.md`, `AGENTS_AND_SKILLS.md`, `docs/architecture.md` (update), `docs/prd.md`, `test/e2e/*.spec.ts`, a git tag
Implementation: document the actual agent(s)/skill(s) used building this (be honest — if this plan itself was agent-assisted, that's your custom-agent story: "spec-to-implementation agent constrained to this document" is a legitimate, defensible entry), write PRD user stories/acceptance criteria for `pmx diagnose` and "Run CI," add 3–4 Playwright specs (dispatch a run, assert state transitions render, assert diagnose output for a known pattern), tag `v0.1.0`.
Definition of done: all 5 non-negotiable checkpoints present and checkable in the repo.

---

## 16. DEMO SCRIPT (4 minutes)

**0:00–0:30 — The problem.** "This pattern should match this file. It doesn't. Here's what every developer does today" — show a generic grep/trial-and-error framing, no slides.

**0:30–1:15 — `pmx diagnose` live in a terminal.** Type `pmx diagnose "src/**/*.ts" "src/index.tsx"` on stage. Point at the verdict, the exact failing literal, the re-verified suggested fix. Run the suggested pattern to show it now matches.

**1:15–1:45 — Same thing, in-browser, via WASM.** Switch to the Foundry's Diagnose tab, type the same inputs, show identical output rendered instantly with no network tab activity — "this is the same compiled Go code, running in your browser."

**1:45–2:45 — Real CI.** Click "Run CI" in the Foundry. Immediately open the GitHub Actions tab in a second window/tab and point at the run appearing there in real time — QUEUED → RUNNING with matching job names. Let it run one job to completion on screen (pick the fastest — `quality`), show PASSED with the real run ID and duration.

**2:45–3:15 — Fuzz evidence.** Click "Run Fuzz." While it runs, explain what `FuzzScan`/`FuzzParse`/`FuzzIsMatch` actually do. Show the exec count and corpus growth numbers landing from the real run.

**3:15–3:45 — Compatibility.** `pmx compat` in terminal, show the real N/50 parity number and, if any exist, the honest divergence list.

**3:45–4:00 — Close.** "Everything you just watched — the diagnosis, the CI run, the fuzz numbers — is real. No mock data, no simulated logs. Here's the run on GitHub if you want to check it yourselves" — paste the Actions run URL in chat/screen.

---

## 17. JUDGE QUESTIONS

**"Why not just use Picomatch JS?"** — Different runtime target: Go tooling (build systems, CLIs, GitHub Actions written in Go) wants a Go-native, zero-dependency matcher rather than shelling out to Node. `pmx` is specifically valuable *because* it's Go-native — no npm install required to debug a glob in a Go-based CI pipeline.

**"Why does this need a CLI?"** — Because the workflow it targets (debugging a pattern before it goes into a config file) is a pre-commit, local, offline action — a web app can't be in a developer's terminal or their CI job. The CLI and the web dashboard are two views of the identical internals, not two products.

**"Why is this Track B?"** — It is literally the "helps software teams debug their own tooling configuration" case the track description gives (bug-triage-adjacent: this is bug triage for glob-driven config).

**"Why can't I use `filepath.Glob`?"** — `filepath.Glob` doesn't support `**` globstars, brace expansion, extglobs, or negation — it's a different, much smaller feature set than picomatch's. That gap is exactly why picomatch (JS) became the de facto standard glob engine across the JS ecosystem, and why porting it to Go has independent value beyond this hackathon.

**"How do you prove compatibility?"** — Point at `compat/expected.json`: real output from real `micromatch/picomatch`, committed, diffed by a real Go test. Report the true fraction, not a claimed 100%.

**"Is the dashboard real?"** — Yes, and provably so: click "Run CI" and watch the run appear on github.com/…/actions in a separate tab within seconds, with a matching run ID.

**"Are these CI results real?"** — Same answer — the run ID and URL are your proof; offer to let a judge dispatch a run themselves.

**"What happens when something fails?"** — Show it: state the failing job/step name and first error line, sourced from the real job logs, not a canned failure message. (Consider deliberately breaking a benchmark or test on a branch before the demo so you have a *real* failure to show, rather than only ever demoing the happy path.)

**"What developer problem are you actually solving?"** — Restate §4 in one sentence: turning a silent boolean match failure into an explained, verified root cause.

---

## 18. 9/10 SCORECARD

| Dimension | Score | Why |
|---|---|---|
| Developer usefulness | 8 | Real, narrow, immediately usable problem; ceiling capped slightly because the addressable moment (debugging a glob) is real but not daily-driver frequent for most devs |
| Technical depth | 9 | Built on a genuinely correct hand-rolled scanner/parser/regex-compiler/matcher, with real fuzzing and real differential compatibility testing behind it |
| Originality | 8 | "Debugger for an existing correct engine" is a less common submission shape than "yet another wrapper app" |
| Demo quality | 9 | Terminal → browser → live GitHub Actions run, in under 4 minutes, with an independently verifiable URL at the end |
| Track alignment | 9 | Matches the track's own example categories almost exactly |
| Feasibility | 8 | Almost everything is a thin layer over code that already exists; the only genuinely new logic is `diagnose`'s divergence-finder and the CI dispatch/poll routes, both scoped to hours not days |
| Production credibility | 8 | Zero-dependency Go core, real CI, real fuzzing, real compat proof — held back one point because a single-PAT auth model is a hackathon-grade, not production-grade, security posture (say so out loud, don't pretend otherwise) |
| Engineering rigor | 9 | Fmt/lint/vet/test/race/fuzz/bench/compat all real and CI-enforced |
| Judge memorability | 9 | "They clicked a button and I watched the real GitHub Actions run happen in another tab" is the kind of specific detail judges remember over generic dashboards |

Nothing here is below 8; no redesign required. The two things that would drag every score down are exactly the two findings in §1 — ship those fixes first.

---

## 19. DO NOT BUILD

- The existing "CI/CD pipeline simulator" / "Trigger Integration Pipeline" `setTimeout` sequence — delete it, don't extend it.
- Any LLM-based "explain why this didn't match" feature — the deterministic explanation is the differentiator; don't dilute it.
- User accounts, saved-pattern history, multi-repo support, teams/orgs.
- A generic CI/workflow engine — you have GitHub Actions, use it, don't rebuild a scheduler.
- Fabricated coverage percentages for the native fuzz targets.
- A live Node.js subprocess in the CLI or browser runtime path for compatibility checking — precompute fixtures instead (§11).
- More than the 6 CLI commands listed in §7. Resist the urge to add `pmx scan`/`pmx parse` as separate top-level commands — `explain` already exposes both; extra commands are just extra demo-script surface area to manage under time pressure.
- Coverage/percentage claims in the README that aren't regenerated by a real, current CI run (replace the static "Apple M2, Go 1.21" sample numbers before demo day).

---

## 20. FINAL ONE-PARAGRAPH PITCH

`pmx` is a debugger for glob patterns, built directly on the scanner, parser, regex compiler, and matcher your team already ported correctly from JavaScript picomatch to Go: given a pattern and a path that don't match, it walks the same internal stages the matcher itself uses, names the exact literal or flag responsible for the mismatch, and proposes a fix it has already re-verified with a real match call — available identically as a local CLI, an in-browser WASM tool, and a GitHub Action — while the accompanying Foundry dashboard proves the engine's correctness isn't just claimed but demonstrated live, by dispatching real GitHub Actions runs for CI, fuzzing, benchmarking, and JS-compatibility checking and streaming back their actual run IDs, job statuses, and logs, with zero simulated states anywhere in the product.