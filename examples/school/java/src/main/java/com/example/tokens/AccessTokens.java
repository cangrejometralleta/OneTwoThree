package com.example.tokens;

import com.example.app.TokenIssuer;
import com.example.codec.Json;
import com.example.faults.Fault;

import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.security.GeneralSecurityException;
import java.security.MessageDigest;
import java.time.Clock;
import java.util.Base64;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.Optional;

import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;

/**
 * AccessTokens Mints and Reads Bearer Tokens signed with HMAC-SHA256, javax.crypto only.
 * A Token is two base64url Parts without padding: the Claims, then the Signature over them.
 * Reference: https://www.rfc-editor.org/rfc/rfc2104
 */
public final class AccessTokens implements TokenIssuer {

  /** TOKEN_IS_INVALID Answers every Token the Guard cannot Trust; a Caller Learns only that. */
  public static final Fault TOKEN_IS_INVALID = Fault.refuseUnprovenCaller("token is Invalid or Expired");

  private static final String SIGNATURE_ALGORITHM = "HmacSHA256";
  private static final Base64.Encoder ENCODER = Base64.getUrlEncoder().withoutPadding();
  private static final Base64.Decoder DECODER = Base64.getUrlDecoder();

  private final byte[] secret;
  private final long lifeSeconds;
  private final Clock clock;

  private AccessTokens(String secret, long lifeSeconds, Clock clock) {
    this.secret = secret.getBytes(StandardCharsets.UTF_8);
    this.lifeSeconds = lifeSeconds;
    this.clock = clock;
  }

  /** buildAccessTokens Casts the Issuer with its Secret, how long a Token Lives, and a Clock. */
  public static AccessTokens buildAccessTokens(String secret, long lifeSeconds, Clock clock) {
    return new AccessTokens(secret, lifeSeconds, clock);
  }

  /** issueAccessToken Mints a Token for one Subject. An empty Secret Signs nothing. */
  @Override
  public String issueAccessToken(String subject) {
    if (secret.length == 0 || subject.isEmpty()) {
      throw TOKEN_IS_INVALID;
    }

    Map<String, Object> claims = new LinkedHashMap<>();
    claims.put("sub", subject);
    claims.put("exp", readNowSeconds() + lifeSeconds);
    String encoded = ENCODER.encodeToString(Json.writeJsonText(claims).getBytes(StandardCharsets.UTF_8));

    return encoded + "." + sign(encoded);
  }

  /**
   * readCaller Names the Subject a Token Proves.
   * The Comparison Reads forward: a Token is Dead when now is after its Deadline.
   */
  @Override
  public String readCaller(String token) {
    String[] parts = token.split("\\.", -1);
    if (secret.length == 0 || parts.length != 2
        || !MessageDigest.isEqual(parts[1].getBytes(StandardCharsets.US_ASCII),
            sign(parts[0]).getBytes(StandardCharsets.US_ASCII))) {
      throw TOKEN_IS_INVALID;
    }

    Map<?, ?> claims = readClaims(parts[0]).orElseThrow(() -> TOKEN_IS_INVALID);
    if (!(claims.get("sub") instanceof String subject) || subject.isEmpty()
        || !(claims.get("exp") instanceof BigInteger deadline)) {
      throw TOKEN_IS_INVALID;
    }
    if (BigInteger.valueOf(readNowSeconds()).compareTo(deadline) > 0) {
      throw TOKEN_IS_INVALID;
    }

    return subject;
  }

  private static Optional<Map<?, ?>> readClaims(String encoded) {
    try {
      return Json.readJsonText(DECODER.decode(encoded))
          .filter(Map.class::isInstance)
          .map(Map.class::cast);
    } catch (IllegalArgumentException unreadable) {
      return Optional.empty();
    }
  }

  /** readNowSeconds Truncates Now to whole unix Seconds. */
  private long readNowSeconds() {
    return clock.instant().getEpochSecond();
  }

  private String sign(String encoded) {
    try {
      Mac mac = Mac.getInstance(SIGNATURE_ALGORITHM);
      mac.init(new SecretKeySpec(secret, SIGNATURE_ALGORITHM));

      return ENCODER.encodeToString(mac.doFinal(encoded.getBytes(StandardCharsets.US_ASCII)));
    } catch (GeneralSecurityException impossible) {
      throw new IllegalStateException("the JDK Ships HmacSHA256", impossible);
    }
  }
}
