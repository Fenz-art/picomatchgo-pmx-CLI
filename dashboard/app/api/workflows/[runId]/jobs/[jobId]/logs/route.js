import { NextResponse } from 'next/server';
import { getJobLogs } from '../../../../../github/utils';

export const dynamic = 'force-dynamic';

export async function GET(_request, { params }) {
  try {
    const { jobId } = await params;
    const logs = await getJobLogs(jobId);
    return new NextResponse(logs, {
      headers: {
        'Content-Type': 'text/plain; charset=utf-8',
      },
    });
  } catch (error) {
    return NextResponse.json({ error: error.message }, { status: 500 });
  }
}
