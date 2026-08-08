import { NextResponse } from 'next/server';
import { getRun, getJobs, getArtifacts } from '../../github/utils';

export const dynamic = 'force-dynamic';

export async function GET(_request, { params }) {
  try {
    const { runId } = await params;
    const run = await getRun(runId);
    const jobs = await getJobs(runId);
    const artifacts = await getArtifacts(runId);

    return NextResponse.json({
      ...run,
      jobs: jobs.jobs || [],
      artifacts: artifacts.artifacts || [],
    });
  } catch (error) {
    return NextResponse.json({ error: error.message }, { status: 500 });
  }
}
