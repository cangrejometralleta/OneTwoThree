import { createHmac, timingSafeEqual } from 'node:crypto';
import { parseJson, stringifyJson } from '../codec/json.ts';
import type { TokenIssuer } from '../app/ports.ts';
const b64 = (value: string | Uint8Array) => Buffer.from(value).toString('base64url');
export class HmacTokenIssuer implements TokenIssuer {
  private secret: string; private lifeSeconds: bigint;
  constructor(secret: string, lifeSeconds: bigint) { this.secret=secret; this.lifeSeconds=lifeSeconds; }
  issue(nowSeconds = BigInt(Math.floor(Date.now() / 1000))): string {
    if (!this.secret) throw new Error('empty secret');
    const payload = b64(stringifyJson({ sub: 'student-registry', exp: nowSeconds + this.lifeSeconds }));
    return `${payload}.${b64(createHmac('sha256', Buffer.from(this.secret, 'utf8')).update(payload, 'ascii').digest())}`;
  }
  verify(token: string, nowSeconds = BigInt(Math.floor(Date.now() / 1000))): string | undefined {
    try {
      if (!this.secret || token.split('.').length !== 2) return undefined;
      const parts = token.split('.');
      const payload = parts[0]!;
      const supplied = parts[1]!;
      if (!/^[A-Za-z0-9_-]+$/.test(payload) || !/^[A-Za-z0-9_-]{43}$/.test(supplied)) return undefined;
      if (Buffer.from(payload,'base64url').toString('base64url') !== payload || Buffer.from(supplied,'base64url').toString('base64url') !== supplied) return undefined;
      const expected = createHmac('sha256', Buffer.from(this.secret, 'utf8')).update(payload, 'ascii').digest();
      const actual = Buffer.from(supplied, 'base64url'); if (actual.length !== expected.length || !timingSafeEqual(actual, expected)) return undefined;
      const text = new TextDecoder('utf-8',{fatal:true}).decode(Buffer.from(payload,'base64url'));
      const claims = parseJson(text) as Record<string, unknown>;
      if (Object.keys(claims).length !== 2 || typeof claims.sub !== 'string' || !claims.sub || typeof claims.exp !== 'bigint' || nowSeconds > claims.exp) return undefined;
      return claims.sub;
    } catch { return undefined; }
  }
}
