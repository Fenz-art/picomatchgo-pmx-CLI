import { getGitHubRepo, getGitHubToken } from './config.mjs';

export const dynamic = 'force-dynamic';

const REPO = getGitHubRepo();
const TOKEN = getGitHubToken();
const API_BASE = process.env.GITHUB_API_URL || 'https://api.github.com';

function getHeaders() {
  const headers = {
    Accept: 'application/vnd.github+json',
  };

  if (!TOKEN) {
    throw new Error(
      'Missing GitHub token environment variable. Set FOUNDRY_GITHUB_TOKEN or GITHUB_TOKEN in your deployment to enable GitHub Actions integration.'
    );
  }

  headers.Authorization = `Bearer ${TOKEN}`;
  headers['X-GitHub-Api-Version'] = '2022-11-28';
  return headers;
}

async function githubFetch(path, init = {}) {
  const response = await fetch(`${API_BASE}${path}`, {
    ...init,
    headers: {
      ...getHeaders(),
      ...(init.headers || {}),
    },
  });

  if (!response.ok) {
    const body = await response.text();
    throw new Error(`GitHub API ${path} failed with ${response.status}: ${body}`);
  }

  const text = await response.text();
  if (!text) {
    return null;
  }

  return JSON.parse(text);
}

export function resolveWorkflowPath(workflow) {
  if (!workflow || workflow === 'production-ci') {
    return 'ci.yml';
  }
  throw new Error('Unsupported workflow. Only production-ci can be dispatched from Foundry.');
}

export async function dispatchWorkflow(workflow = 'production-ci', ref = 'main', inputs = {}) {
  const workflowPath = resolveWorkflowPath(workflow);

  await githubFetch(`/repos/${REPO}/actions/workflows/${workflowPath}/dispatches`, {
    method: 'POST',
    body: JSON.stringify({ ref, inputs }),
  });

  for (let attempt = 0; attempt < 10; attempt += 1) {
    const runs = await githubFetch(`/repos/${REPO}/actions/workflows/${workflowPath}/runs?event=workflow_dispatch&per_page=5&branch=${encodeURIComponent(ref)}`);
    if (Array.isArray(runs.workflow_runs) && runs.workflow_runs.length > 0) {
      return runs.workflow_runs[0];
    }
    await new Promise((resolve) => setTimeout(resolve, 1200));
  }

  throw new Error('Unable to resolve workflow run after dispatch.');
}

export async function listWorkflowRuns(workflow = 'production-ci') {
  const workflowPath = resolveWorkflowPath(workflow);
  const result = await githubFetch(`/repos/${REPO}/actions/workflows/${workflowPath}/runs?per_page=8&branch=main`);
  return result.workflow_runs || [];
}

export async function getRun(runId) {
  return githubFetch(`/repos/${REPO}/actions/runs/${runId}`);
}

export async function getJobs(runId) {
  return githubFetch(`/repos/${REPO}/actions/runs/${runId}/jobs?per_page=50`);
}

export async function getArtifacts(runId) {
  return githubFetch(`/repos/${REPO}/actions/runs/${runId}/artifacts?per_page=50`);
}

export async function getJobDetails(jobId) {
  return githubFetch(`/repos/${REPO}/actions/jobs/${jobId}`);
}

export async function getJobLogs(jobId) {
  const response = await fetch(`${API_BASE}/repos/${REPO}/actions/jobs/${jobId}/logs`, {
    headers: getHeaders(),
  });

  if (!response.ok) {
    const body = await response.text();
    throw new Error(`GitHub job logs failed with ${response.status}: ${body}`);
  }

  return response.text();
}
