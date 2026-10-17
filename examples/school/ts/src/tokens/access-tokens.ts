import { createHmac, timingSafeEqual } from 'node:crypto';

import type { TokenIssuer } from '../app/ports.ts';
import { readJsonText, writeJsonText } from '../codec/json.ts';
import { refuseUnprovenCaller } from '../faults/faults.ts';

/** TOKEN_IS_INVALID Answers every Token the Guard cannot Trust; a Caller Learns only that. */
export const TOKEN_IS_INVALID = refuseUnprovenCaller('token is Invalid or Expired');

const MILLISECONDS_PER_SECOND = 1000;

/** Clock Tells Now, in Milliseconds since the Epoch. */
export type Clock = () => number;

/**
 * AccessTokens Mints and Reads Bearer Tokens signed with HMAC-SHA256, node:crypto only.
 * A Token is two base64url Parts without padding: the Claims, then the Signature over them.
 * Reference: https://www.rfc-editor.org/rfc/rfc2104
 */
export class AccessTokens implements TokenIssuer {
  readonly #secret: Buffer;
  readonly #lifeSeconds: bigint;
  readonly #clock: Clock;

  constructor(secret: string, lifeSeconds: bigint, clock: Clock = Date.now) {
    this.#secret = Buffer.from(secret, 'utf8');
    this.#lifeSeconds = lifeSeconds;
    this.#clock = clock;
  }

  /** issueAccessToken Mints a Token for one Subject. An empty Secret Signs nothing. */
  issueAccessToken(subject: string): string {
    if (this.#secret.length === 0 || subject === '') {
      throw TOKEN_IS_INVALID;
    }

    const claims = writeJsonText({ sub: subject, exp: this.#readNowSeconds() + this.#lifeSeconds });
    const encoded = Buffer.from(claims, 'utf8').toString('base64url');

    return `${encoded}.${this.#sign(encoded)}`;
  }

  /** readCaller Names the Subject a Token Proves. The Comparison Reads forward: a Token is Dead when now is after its Deadline. */
  readCaller(token: string): string {
    const parts = token.split('.');
    const [encoded = '', signature = ''] = parts;
    const expected = Buffer.from(this.#secret.length === 0 ? '' : this.#sign(encoded), 'ascii');
    const given = Buffer.from(signature, 'ascii');
    if (this.#secret.length === 0 || parts.length !== 2 || given.length !== expected.length || !timingSafeEqual(given, expected)) {
      throw TOKEN_IS_INVALID;
    }

    const claims = readJsonText(Buffer.from(encoded, 'base64url'));
    const subject = claims instanceof Map ? claims.get('sub') : undefined;
    const deadline = claims instanceof Map ? claims.get('exp') : undefined;
    if (typeof subject !== 'string' || subject === '' || typeof deadline !== 'bigint') {
      throw TOKEN_IS_INVALID;
    }
    if (this.#readNowSeconds() > deadline) {
      throw TOKEN_IS_INVALID;
    }

    return subject;
  }

  /** readNowSeconds Truncates Now to whole unix Seconds. */
  #readNowSeconds(): bigint {
    return BigInt(Math.floor(this.#clock() / MILLISECONDS_PER_SECOND));
  }

  #sign(encoded: string): string {
    return createHmac('sha256', this.#secret).update(encoded, 'ascii').digest('base64url');
  }
}
