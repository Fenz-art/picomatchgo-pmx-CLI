import { NextResponse } from 'next/server';
import { execFile } from 'node:child_process';
import { promisify } from 'node:util';
import path from 'node:path';

const execFileAsync = promisify(execFile);
const repoRoot = path.resolve(process.cwd(), '..');

const validationRegistry = {
  help: { command: 'pmx help', label: 'pmx help', args: ['help'], mode: 'local', category: 'CLI', executor: 'LocalExecutor' },
  match: { command: 'pmx match', label: 'pmx match', args: ['match', '*.js', 'foo.js'], mode: 'local', category: 'CLI', executor: 'LocalExecutor' },
  scan: { command: 'pmx scan', label: 'pmx scan', args: ['scan', '**/parser/*.go'], mode: 'local', category: 'CLI', executor: 'LocalExecutor' },
  parse: { command: 'pmx parse', label: 'pmx parse', args: ['parse', 'foo/{bar,baz}/@(a|b).go'], mode: 'local', category: 'CLI', executor: 'LocalExecutor' },
  explain: { command: 'pmx explain', label: 'pmx explain', args: ['explain', '**/*.go', '--input', 'src/parser/scan.go'], mode: 'local', category: 'CLI', executor: 'LocalExecutor' },
  validate: { command: 'pmx validate', label: 'pmx validate', args: ['validate', '*.go', '--input', 'foo.go'], mode: 'local', category: 'CLI', executor: 'LocalExecutor' },
  compat: { command: 'pmx compat', label: 'pmx compat', args: ['compat', '--suite', 'basic'], mode: 'local', category: 'Compatibility', executor: 'LocalExecutor' },
  regression: { command: 'pmx regression --json', label: 'pmx regression --json', args: ['regression', '--json'], mode: 'local', category: 'Regression', executor: 'LocalExecutor' },
  bench: { command: 'pmx bench', label: 'pmx bench', args: ['bench'], mode: 'local', category: 'Reliability', executor: 'LocalExecutor' },
  fuzz: { command: 'pmx fuzz --json --time 5s', label: 'pmx fuzz --json --time 5s', args: ['fuzz', '--json', '--time', '5s'], mode: 'local', category: 'Reliability', executor: 'LocalExecutor' },
  doctor: { command: 'pmx doctor', label: 'pmx doctor', args: ['doctor'], mode: 'local', category: 'Doctor', executor: 'LocalExecutor' },
  doctorJson: { command: 'pmx doctor --json', label: 'pmx doctor --json', args: ['doctor', '--json'], mode: 'local', category: 'Doctor', executor: 'LocalExecutor' },
  doctorCi: { command: 'pmx doctor --ci', label: 'pmx doctor --ci', args: ['doctor', '--ci'], mode: 'local', category: 'Doctor', executor: 'LocalExecutor' },
  agentInspect: { command: 'pmx agent inspect --json', label: 'pmx agent inspect --json', args: ['agent', 'inspect', '--json'], mode: 'local', category: 'ADLC', executor: 'LocalExecutor' },
  agentCheck: { command: 'pmx agent check --json', label: 'pmx agent check --json', args: ['agent', 'check', '--json'], mode: 'local', category: 'ADLC', executor: 'LocalExecutor' },
  ci: { command: 'pmx ci --json', label: 'pmx ci --json', args: ['ci', '--json'], mode: 'local', category: 'CI', executor: 'LocalExecutor' },
};

function makeSummary(stdout, stderr, exitCode) {
  const output = `${stdout || ''}${stderr || ''}`.trim();
  if (!output) return exitCode === 0 ? 'Command completed successfully.' : 'Command returned no output.';

  const lines = output.split(/\r?\n/).map((line) => line.trim()).filter(Boolean);
  const lastLine = lines.at(-1) || output;
  return lastLine.slice(0, 180);
}

function classifyStatus(exitCode, output) {
  if (exitCode === 0) {
    const normalized = (output || '').toLowerCase();
    if (normalized.includes('warning') || normalized.includes('result: warning')) {
      return 'warn';
    }
    return 'pass';
  }
  return 'fail';
}

function normalizeResult(commandId, spec, exitCode, durationMs, stdout, stderr) {
  const output = `${stdout || ''}${stderr || ''}`.trim();
  const status = classifyStatus(exitCode, output);
  let parsed = null;
  try {
    parsed = stdout && /^[\s\r\n]*[\[{]/.test(stdout) ? JSON.parse(stdout) : null;
  } catch {
    // Human-readable commands intentionally have no parsed payload.
  }
  return {
    id: commandId,
    command: spec.command,
    name: spec.label,
    category: spec.category,
    mode: spec.mode,
    executor: spec.executor,
    ok: status !== 'fail',
    status,
    exitCode,
    durationMs,
    stdout: stdout || '',
    stderr: stderr || '',
    output,
    summary: makeSummary(stdout, stderr, exitCode),
    parsed,
  };
}

export function buildValidationRun(results, meta = {}) {
  const safeResults = Array.isArray(results) ? results : [];
  const startedAt = meta.startedAt || new Date().toISOString();
  const finishedAt = meta.finishedAt || new Date().toISOString();
  const passed = safeResults.filter((entry) => entry?.status === 'pass').length;
  const failed = safeResults.filter((entry) => entry?.status === 'fail').length;
  const warnings = safeResults.filter((entry) => entry?.status === 'warn').length;
  const status = failed > 0 ? 'fail' : warnings > 0 ? 'warn' : safeResults.length > 0 ? 'pass' : 'fail';

  return {
    id: meta.id || `validation-${Date.now()}`,
    type: meta.type || 'suite',
    status,
    passed,
    failed,
    warnings,
    total: safeResults.length,
    startedAt,
    finishedAt,
    stages: safeResults.map((entry) => ({
      name: entry?.id || 'unknown',
      label: entry?.name || entry?.command || 'Unknown validation',
      category: entry?.category || 'CLI',
      executor: entry?.executor || 'LocalExecutor',
      status: entry?.status || 'fail',
      exitCode: typeof entry?.exitCode === 'number' ? entry.exitCode : 1,
      durationMs: typeof entry?.durationMs === 'number' ? entry.durationMs : 0,
      stdout: entry?.stdout || '',
      stderr: entry?.stderr || '',
      output: entry?.output || '',
      diagnostics: entry?.status === 'fail'
        ? [{ level: 'error', message: entry?.summary || 'Command failed.' }]
        : entry?.status === 'warn'
          ? [{ level: 'warning', message: entry?.summary || 'Command reported warnings.' }]
          : [],
      parsed: entry?.parsed || null,
    })),
    result: {
      status,
      passed,
      failed,
      warnings,
      total: safeResults.length,
    },
  };
}

function buildCommandEnv(spec) {
  const env = { ...process.env };
  if (spec?.args?.[0] === 'agent' && spec.args[1] === 'check' || spec?.args?.[0] === 'ci') {
    env.PMX_AGENT_CHECK_SKIP_CI = '1';
  }
  return env;
}

async function runValidationCommand(commandId, spec) {
  const startedAt = Date.now();

  try {
    const { stdout, stderr } = await execFileAsync('go', ['run', './cmd/pmx', ...spec.args], {
      cwd: repoRoot,
      timeout: 300000,
      maxBuffer: 1024 * 1024,
      env: buildCommandEnv(spec),
    });

    return normalizeResult(commandId, spec, 0, Date.now() - startedAt, stdout, stderr);
  } catch (error) {
    const stdout = error.stdout ? String(error.stdout) : '';
    const stderr = error.stderr ? String(error.stderr) : '';
    const exitCode = typeof error.code === 'number' ? error.code : 1;
    const output = `${stdout || ''}${stderr || ''}`.trim() || error.message || 'Command failed.';
    return normalizeResult(commandId, spec, exitCode, Date.now() - startedAt, stdout, stderr || output);
  }
}

export async function GET() {
  return NextResponse.json({
    registry: Object.entries(validationRegistry).map(([id, item]) => ({ id, ...item })),
    repoRoot,
  });
}

export async function POST(request) {
  try {
    const body = await request.json();
    const requested = Array.isArray(body?.commands)
      ? body.commands
      : body?.command
        ? [body.command]
        : Object.keys(validationRegistry);

    const startedAt = new Date().toISOString();
    const results = [];
    let passed = 0;
    let failed = 0;

    for (const commandId of requested) {
      const spec = validationRegistry[commandId];
      if (!spec) {
        const result = {
          id: commandId,
          command: commandId,
          name: commandId,
          category: 'Unknown',
          mode: 'local',
          executor: 'LocalExecutor',
          ok: false,
          status: 'fail',
          exitCode: 1,
          durationMs: 0,
          stdout: '',
          stderr: `Unknown validation command: ${commandId}`,
          output: `Unknown validation command: ${commandId}`,
          summary: `Unknown validation command: ${commandId}`,
        };
        results.push(result);
        failed += 1;
        continue;
      }

      const result = await runValidationCommand(commandId, spec);
      results.push(result);
      if (result.status === 'pass') {
        passed += 1;
      } else if (result.status === 'fail') {
        failed += 1;
      }
    }

    const ok = failed === 0;
    const status = ok ? 'pass' : (passed > 0 ? 'warn' : 'fail');
    const finishedAt = new Date().toISOString();
    const run = buildValidationRun(results, {
      id: `validation-${Date.now()}`,
      type: 'complete',
      startedAt,
      finishedAt,
    });

    return NextResponse.json({
      ok,
      status,
      passed,
      failed,
      total: results.length,
      results,
      run,
    });
  } catch (error) {
    return NextResponse.json({
      ok: false,
      status: 'fail',
      error: error.message,
      results: [],
    }, { status: 500 });
  }
}
