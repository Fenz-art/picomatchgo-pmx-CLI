# PMX Agents and Skills

## Custom agent: PMX Validation Agent

Definition: `.codex/agents/pmx-validation-agent.md`

The agent performs a bounded inspect → diagnose → fix → check → regression workflow for PMX projects. It accepts a repository path and a task description, consumes `pmx agent inspect --json`, `pmx doctor --json`, and `pmx agent check --json`, and returns structured evidence: changed files, check results, failures, and next actions. Invoke it for validation failures, submission readiness, or a Foundry/CLI execution discrepancy.

## Custom skill: PMX Doctor and ADLC Validation

Definition: `.codex/skills/pmx-doctor-adlc/SKILL.md`

The skill defines the repeatable PMX diagnostic workflow for JSON-based doctor and ADLC validations. It executes `pmx doctor --json`, `pmx agent inspect --json`, and `pmx agent check --json`, then uses the command result payloads to drive repair, regression, and final validation reporting.

This skill defines the repeatable PMX diagnostic workflow. It uses doctor and ADLC JSON contracts rather than scraping prose, distinguishes warnings from failures, executes the relevant regression command after a repair, and reports the actual exit status and captured evidence.

## Lifecycle

```text
INSPECT -> DIAGNOSE -> FIX -> CHECK -> REGRESSION -> CI -> REPORT
```

The custom agent invokes the skill at the diagnose/check stages. Foundry may display the same command evidence through `/api/validation`; GitHub workflow state is obtained only through the server-side GitHub API routes.
