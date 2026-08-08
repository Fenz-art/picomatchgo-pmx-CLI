import test from 'node:test';
import assert from 'node:assert/strict';
import { buildValidationRun } from './route.js';

test('buildValidationRun creates a structured run with per-stage evidence', () => {
  const run = buildValidationRun([
    {
      id: 'help',
      command: 'pmx help',
      name: 'pmx help',
      category: 'CLI',
      mode: 'local',
      executor: 'LocalExecutor',
      ok: true,
      status: 'pass',
      exitCode: 0,
      durationMs: 42,
      stdout: 'pmx — Picomatch Go Developer Diagnostics',
      stderr: '',
      output: 'pmx — Picomatch Go Developer Diagnostics',
      summary: 'pmx — Picomatch Go Developer Diagnostics',
    },
    {
      id: 'agentCheck',
      command: 'pmx agent check --json',
      name: 'pmx agent check --json',
      category: 'ADLC',
      mode: 'local',
      executor: 'LocalExecutor',
      ok: true,
      status: 'pass',
      exitCode: 0,
      durationMs: 64,
      stdout: '{"status":"pass"}',
      stderr: '',
      output: '{"status":"pass"}',
      summary: '{"status":"pass"}',
    },
  ], { type: 'complete' });

  assert.equal(run.type, 'complete');
  assert.equal(run.status, 'pass');
  assert.equal(run.passed, 2);
  assert.equal(run.failed, 0);
  assert.equal(run.stages.length, 2);
  assert.equal(run.stages[0].name, 'help');
  assert.equal(run.stages[0].executor, 'LocalExecutor');
  assert.ok(run.id);
});
