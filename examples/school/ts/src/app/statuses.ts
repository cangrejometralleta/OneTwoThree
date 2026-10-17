import { readFaultKind, type Kind } from '../faults/faults.ts';

/**
 * FAULT_STATUSES Translates each Kind into the Status HTTP Knows it by.
 * The Core Names the Kind; only the Application Layer Knows the Protocol.
 * Reference: https://www.rfc-editor.org/rfc/rfc9110#section-15
 */
const FAULT_STATUSES: ReadonlyMap<Kind, number> = new Map([
  ['invalid-input', 400],
  ['unproven-caller', 401],
  ['missing-record', 404],
  ['taken-value', 409],
  ['unprocessable-value', 422],
  ['unavailable-provider', 503],
]);

/** HTTP_INTERNAL_ERROR Answers a Failure that is Ours. */
export const HTTP_INTERNAL_ERROR = 500;

/** readFaultStatus is the one Place that Turns a Failure into a Number. A Kind without a Status Answers five hundred. */
export function readFaultStatus(failure: unknown): number {
  return FAULT_STATUSES.get(readFaultKind(failure)) ?? HTTP_INTERNAL_ERROR;
}
