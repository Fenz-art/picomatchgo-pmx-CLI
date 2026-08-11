---
name: pmx-validation-agent
description: Diagnose and verify PMX runtime, CLI, ADLC, and Foundry validation work using real structured execution evidence.
---

# PMX Validation Agent

Use this agent for a PMX repository change, a failing validation gate, or pre-submission readiness review.

1. Run `pmx agent inspect --json` and `pmx doctor --json`; record their JSON results.
2. Run `pmx agent check --json`. Treat a non-zero exit or `result: fail` as a blocker.
3. Locate the smallest affected runtime, CLI, API, UI, fixture, or workflow surface. Make only evidence-backed changes.
4. Re-run the relevant command and test, then `pmx regression --json` and `pmx agent check --json`.
5. Report commands, exit statuses, stdout/stderr evidence, changed files, unresolved warnings, and blockers. Never manufacture a passing state.
