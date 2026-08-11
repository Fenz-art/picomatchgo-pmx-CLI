# PMX Agent Integration

PMX supports a lightweight agent-driven lifecycle for coding agents and automation tools.

## Goals

- Make project inspection deterministic and machine-readable.
- Let an agent reason over structured JSON instead of terminal prose.
- Give CI and human developers the same lifecycle contract.

## Lifecycle

The current PMX agent contract covers two essential phases:

1. Inspect the repository and compute an engineering snapshot.
2. Check the repo against validation and diagnostics rules.

## Commands

### Inspect

```bash
pmx agent inspect --json
```

Returns a structured report with:

- project metadata
- diagnostic entries
- overall status

Example shape:

```json
{
  "version": "1",
  "project": {
    "name": "project",
    "ecosystem": "javascript/typescript",
    "package_manager": "pnpm",
    "typescript": true,
    "framework": "next"
  },
  "diagnostics": [
    {
      "id": "PMX-TS-001",
      "severity": "warn",
      "category": "typescript",
      "title": "TypeScript strict mode is disabled",
      "file": "tsconfig.json",
      "message": "TypeScript configuration detected but strict mode is disabled.",
      "suggestion": "Review whether strict mode should be enabled for this project."
    }
  ]
}
```

### Check

```bash
pmx agent check --json
```

Returns the strict ADLC gate payload that agents should treat as the canonical verification contract:

- version
- result
- diagnostics
- checks
- next actions

Example shape:

```json
{
  "version": "1",
  "result": "warn",
  "diagnostics": [
    {
      "code": "PMX-TS-001",
      "severity": "warning",
      "file": "tsconfig.json",
      "message": "TypeScript strict mode is disabled",
      "suggestion": "Review whether strict mode should be enabled for this project."
    }
  ],
  "checks": [
    {"name": "doctor", "status": "warn"},
    {"name": "validate", "status": "pass"},
    {"name": "compat", "status": "pass"}
  ],
  "next_actions": [
    {"action": "fix", "code": "PMX-TS-001"}
  ]
}
```

## Agent Workflow

A coding agent can follow this lifecycle:

```text
pmx agent inspect --json
  -> determine project shape and diagnostics
pmx agent check --json
  -> validate before/after a set of changes
```

This is the minimal ADLC slice for PMX: inspect, diagnose, and gate verification for agent-driven repair loops.

## Integration Guidance

When a tool or agent is integrating with PMX, prefer:

- `--json` output for machine consumption
- stable diagnostic IDs for inspection
- status-based exit behavior for validation gates
- structured next actions rather than terminal parsing

## Current scope

This version intentionally implements the first viable agent contract for PMX:

- `pmx agent inspect --json`
- `pmx agent check --json`

This gives agents a deterministic interface without forcing them to parse human-readable diagnostics.
