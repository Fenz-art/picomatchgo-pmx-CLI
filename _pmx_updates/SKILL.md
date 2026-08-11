---
name: pmx-doctor-adlc
description: Run PMX diagnostics and the ADLC gate through their stable JSON contracts.
---

# PMX Doctor and ADLC Validation

Run from the repository root.

1. Execute `pmx doctor --json [project-dir]` to discover project evidence and diagnostics.
2. Execute `pmx agent inspect --json` for the machine-readable project snapshot.
3. Execute `pmx agent check --json` for the canonical gate payload: `version`, `result`, `diagnostics`, `checks`, and `next_actions`.
4. When a check fails, use its captured detail and diagnostics to choose the smallest repair. Do not infer success from the UI or from command availability.
5. After repair, run the targeted test plus `pmx regression --json`; run `pmx agent check --json` again and report the before/after result and exit codes.

Warnings are evidence, not passes. Escalate failures and unavailable external services (such as GitHub credentials) as explicit blockers.
