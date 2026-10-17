/**
 * A Kind Names what Went wrong, in the Business' own Words.
 * It Names no Protocol: app Translates a Kind into a Status.
 */
export type Kind =
  | 'unexpected'
  | 'invalid-input'
  | 'unproven-caller'
  | 'missing-record'
  | 'taken-value'
  | 'unprocessable-value'
  | 'unavailable-provider';

/** INTERNAL_REASON is all a Caller Learns of a Failure that is Ours. */
export const INTERNAL_REASON = 'internal error';

/**
 * A Fault is a Failure the Program Expected, Carrying the Kind it Deserves.
 * A Failure that Reaches the Edge without one is Unexpected.
 */
export class Fault extends Error {
  readonly kind: Kind;

  constructor(kind: Kind, reason: string) {
    super(reason);
    this.kind = kind;
  }

  get reason(): string {
    return this.message;
  }

  /** sameAs Compares two Faults by what they Declare, not by Identity. */
  sameAs(other: unknown): boolean {
    return other instanceof Fault && other.kind === this.kind && other.reason === this.reason;
  }
}

/** refuseInvalidInput Names a Caller that Sent something the Rules Reject. */
export const refuseInvalidInput = (reason: string): Fault => new Fault('invalid-input', reason);

/** refuseUnprovenCaller Names a Caller the Program cannot Recognise. */
export const refuseUnprovenCaller = (reason: string): Fault => new Fault('unproven-caller', reason);

/** reportMissingRecord Names something the Caller Asked for and We Lack. */
export const reportMissingRecord = (reason: string): Fault => new Fault('missing-record', reason);

/** reportTakenValue Names a Collision with something already Stored. */
export const reportTakenValue = (reason: string): Fault => new Fault('taken-value', reason);

/** refuseUnprocessableValue Names a well Formed Value the Authority that Knows still Denies. */
export const refuseUnprocessableValue = (reason: string): Fault => new Fault('unprocessable-value', reason);

/** reportUnavailableProvider Names a Provider that cannot Answer now; an Unknown is never a Yes. */
export const reportUnavailableProvider = (reason: string): Fault => new Fault('unavailable-provider', reason);

/** findFault Finds the Fault inside any Wrapping of causes. */
function findFault(failure: unknown): Fault | undefined {
  for (let cause = failure; cause !== undefined && cause !== null; cause = (cause as { cause?: unknown }).cause) {
    if (cause instanceof Fault) {
      return cause;
    }
    if (typeof cause !== 'object') {
      return undefined;
    }
  }

  return undefined;
}

/** readFaultKind Finds the Fault's Kind; everything else is Unexpected. */
export function readFaultKind(failure: unknown): Kind {
  return findFault(failure)?.kind ?? 'unexpected';
}

/** readFaultReason Finds the Reason a Fault Declares, and Says nothing of a driver. */
export function readFaultReason(failure: unknown): string {
  return findFault(failure)?.reason ?? INTERNAL_REASON;
}
