import { NextResponse } from 'next/server';
import { getJobs } from '../../../github/utils';

export const dynamic = 'force-dynamic';

export async function GET(_request, { params }) {
  try {
    const { runId } = await params;
    const jobs = await getJobs(runId);
    return NextResponse.json({ jobs: jobs.jobs || [] });
  } catch (error) {
    return NextResponse.json({ error: error.message }, { status: 500 });
  }
}
