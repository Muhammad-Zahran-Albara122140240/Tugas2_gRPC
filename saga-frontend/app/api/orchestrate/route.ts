import { NextRequest, NextResponse } from 'next/server';
import { exec } from 'child_process';
import { promisify } from 'util';

const execPromise = promisify(exec);

export async function POST(req: NextRequest) {
  const body = await req.json();
  const { item, amount, address, simulateShippingFailure } = body;

  const simulateFlag = simulateShippingFailure ? "--fail-shipping" : "";
  const command = `go run ../backend/orchestrator/server.go --item="${item}" --amount=${amount} --address="${address}" ${simulateFlag}`;

  try {
    const { stdout, stderr } = await execPromise(command);

    const output = `${stdout}\n${stderr}`.trim(); // Gabung dan trim
    return NextResponse.json({ output });
  } catch (error: any) {
    const output = `${error.stdout || ''}\n${error.stderr || ''}`.trim();
    console.error('Exec error:', error.message);
    return NextResponse.json({ error: error.message, output }, { status: 500 });
  }
}
