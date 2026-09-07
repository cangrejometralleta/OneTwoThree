package com.example.tokens;

import com.example.faults.Fault;
import com.example.school.TokenIssuer;

import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;

import java.nio.charset.StandardCharsets;
import java.time.Instant;
import java.util.Base64;
import java.util.function.Supplier;

/** AccessTokens Fulfils TokenIssuer with the standard Library only. */
public final class AccessTokens implements TokenIssuer {

  /**
   * ErrTokenIsInvalid Covers every Way a Token can Fail.
   * A Caller who Cannot Prove itself Gets one Answer, never a Reason why.
   */
  public static final Fault TOKEN_IS_INVALID = Fault.refuseUnprovenCaller("token is Invalid or Expired");

  private final byte[] secret;
  private final long lifeSeconds;
  private final Supplier<Instant> now;

  private AccessTokens(String secret, long lifeSeconds, Supplier<Instant> now) {
    this.secret = secret.getBytes(StandardCharsets.UTF_8);
    this.lifeSeconds = lifeSeconds;
    this.now = now;
  }

  /** BuildAccessTokens Takes the two Values Validated at Startup. */
  public static AccessTokens buildAccessTokens(String secret, int lifeSeconds) {
    return new AccessTokens(secret, lifeSeconds, Instant::now);
  }

  /** BuildAccessTokensReading Hands a Test a Clock it can Move. */
  public static AccessTokens buildAccessTokensReading(String secret, int lifeSeconds, Supplier<Instant> clock) {
    return new AccessTokens(secret, lifeSeconds, clock);
  }

  /** IssueAccessToken Signs a Subject together with its Deadline. */
  @Override
  public String issueAccessToken(String subject) {
    if (secret.length == 0) {
      throw TOKEN_IS_INVALID;
    }

    String claim = subject + "." + now.get().plusSeconds(lifeSeconds).getEpochSecond();

    return claim + "." + signTokenClaim(claim);
  }

  /**
   * ReadAccessToken Returns the Subject, or Says why it Cannot.
   * The Comparison Reads forward: a Token Dies once its Deadline Passes.
   */
  @Override
  public String readAccessToken(String token) {
    String claim = openSignedClaim(token);
    int dot = claim.lastIndexOf('.');

    long deadline;
    try {
      deadline = Long.parseLong(claim.substring(dot + 1));
    } catch (RuntimeException broken) {
      throw TOKEN_IS_INVALID;
    }

    if (now.get().getEpochSecond() > deadline) {
      throw TOKEN_IS_INVALID;
    }

    return claim.substring(0, dot);
  }

  /** OpenSignedClaim Returns the Claim only when the Signature Holds. */
  private String openSignedClaim(String token) {
    int dot = token == null ? -1 : token.lastIndexOf('.');
    if (dot < 0) {
      throw TOKEN_IS_INVALID;
    }

    String claim = token.substring(0, dot);
    if (!java.security.MessageDigest.isEqual(signTokenClaim(claim).getBytes(StandardCharsets.UTF_8),
        token.substring(dot + 1).getBytes(StandardCharsets.UTF_8))) {
      throw TOKEN_IS_INVALID;
    }

    return claim;
  }

  /** SignTokenClaim Reduces a Claim to its Signature. */
  private String signTokenClaim(String claim) {
    try {
      Mac mac = Mac.getInstance("HmacSHA256");
      mac.init(new SecretKeySpec(secret, "HmacSHA256"));

      return Base64.getUrlEncoder().withoutPadding().encodeToString(mac.doFinal(claim.getBytes(StandardCharsets.UTF_8)));
    } catch (java.security.GeneralSecurityException impossible) {
      throw new IllegalStateException(impossible);
    }
  }
}
