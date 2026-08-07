import { NextResponse } from 'next/server';
import { getArtifacts } from '../../../github/utils';

export const dynamic = 'force-dynamic';

export async function GET(_request, { params }) {
  try {
    const artifacts = await getArtifacts(params.runId);
    return NextResponse.json({ artifacts: artifacts.artifacts || [] });
  } catch (error) {
    return NextResponse.json({ error: error.message }, { status: 500 });
  }
}
