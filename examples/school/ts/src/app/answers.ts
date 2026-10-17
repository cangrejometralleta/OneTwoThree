import { readFaultReason } from '../faults/faults.ts';
import type { Handler, Request } from '../transport/transport.ts';
import type { FailureView } from '../wire/wire.ts';
import type { CallerReader } from './ports.ts';
import { readFaultStatus } from './statuses.ts';
import { readBearerToken } from './validation.ts';

/** Telling is a Script: it Answers with a Value or Fails, and builds no Reply. */
export type Telling = (request: Request) => Promise<unknown>;

/** answerWith is the one Crossing from a Telling to a Handler. The Route Declares the happy Status; a Fault Declares its own. */
export function answerWith(status: number, tell: Telling): Handler {
  return async (request) => {
    try {
      return { status, body: await tell(request) };
    } catch (failure) {
      const view: FailureView = { error: readFaultReason(failure) };

      return { status: readFaultStatus(failure), body: view };
    }
  };
}

/** requireProvenCaller Names the Caller before the Story starts. A Script behind it Reads the Caller and Trusts it. */
export function requireProvenCaller(reader: CallerReader, tell: Telling): Telling {
  return async (request) => tell({ ...request, caller: reader.readCaller(readBearerToken(request)) });
}
