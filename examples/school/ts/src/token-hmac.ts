import { createHmac, timingSafeEqual } from "node:crypto";

import { ErrTokenIsInvalid } from "./domain.js";
import type { TokenIssuer } from "./providers.js";

// HmacTokens Fulfils TokenIssuer with node:crypto only.
// Reference: https://nodejs.org/api/crypto.html#cryptocreatehmacalgorithm-key-options
export class HmacTokens implements TokenIssuer {
  constructor(
    private readonly secret: string,
    private readonly lifeMs: number,
    private readonly now: () => number = Date.now,
  ) {}

  // issueAccessToken Signs a Subject together with its Deadline.
  issueAccessToken(subject: string): string {
    const claim = `${subject}.${this.now() + this.lifeMs}`;

    return `${claim}.${this.signTokenClaim(claim)}`;
  }

  // readAccessToken Returns the Subject, or Throws why it Cannot.
  // The Comparison Reads forward: a Token Dies once its Deadline Passes.
  readAccessToken(token: string): string {
    const claim = this.openSignedClaim(token);
    const [subject, deadline] = claim.split(".");

    if (this.now() > Number(deadline)) throw ErrTokenIsInvalid;

    return subject!;
  }

  // openSignedClaim Returns the Claim only when the Signature Holds.
  private openSignedClaim(token: string): string {
    const cut = token.lastIndexOf(".");
    if (cut < 0) throw ErrTokenIsInvalid;

    const claim = token.slice(0, cut);
    const given = Buffer.from(token.slice(cut + 1));
    const wanted = Buffer.from(this.signTokenClaim(claim));

    if (given.length !== wanted.length || !timingSafeEqual(given, wanted)) throw ErrTokenIsInvalid;

    return claim;
  }

  // signTokenClaim Reduces a Claim to its Signature.
  private signTokenClaim(claim: string): string {
    return createHmac("sha256", this.secret).update(claim).digest("base64url");
  }
}
