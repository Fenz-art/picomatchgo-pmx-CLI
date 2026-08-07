import { NextResponse } from 'next/server';
import { dispatchWorkflow, listWorkflowRuns } from '../github/utils';

export const dynamic = 'force-dynamic';

export async function GET() {
  try {
    const runs = await listWorkflowRuns('production-ci');

    return NextResponse.json({
      workflows: [
        {
          name: 'production-ci',
          path: '.github/workflows/ci.yml',
          runs,
        },
      ],
    });
  } catch (error) {
    return NextResponse.json({ error: error.message }, { status: 500 });
  }
}

export async function POST(request) {
  try {
    const body = await request.json();
    const workflow = body?.workflow || 'production-ci';
    const ref = body?.ref || 'main';
    const inputs = body?.inputs || {};

    const run = await dispatchWorkflow(workflow, ref, inputs);

    return NextResponse.json({
      run,
      workflow,
      message: 'workflow dispatched',
    });
  } catch (error) {
    return NextResponse.json({ error: error.message }, { status: 500 });
  }
}
