const DEFAULT_REPO = 'Fenz-art/picomatchgo-pmx-CLI';

export function resolveRepositorySlug(value) {
  if (!value) {
    return DEFAULT_REPO;
  }

  const repo = String(value).trim().replace(/^\/+/, '');

  if (/^https?:\/\//i.test(repo) || repo.includes('vercel.app') || repo.includes('://')) {
    throw new Error(
      'GitHub repository must be a repository slug like "owner/repo", not a Vercel URL or full GitHub URL.'
    );
  }

  if (!/^[A-Za-z0-9_.-]+\/[A-Za-z0-9_.-]+$/.test(repo)) {
    throw new Error(
      'GitHub repository must be set as "owner/repo" (for example: "kumar-kushang/picomatch-go").'
    );
  }

  return repo;
}

export function getGitHubRepo() {
  const configuredRepo = process.env.FOUNDRY_GITHUB_REPOSITORY || process.env.GITHUB_REPOSITORY;
  return resolveRepositorySlug(configuredRepo);
}

export function getGitHubToken() {
  return process.env.FOUNDRY_GITHUB_TOKEN || process.env.GITHUB_TOKEN || null;
}
