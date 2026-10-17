package com.example.tokens;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

import com.example.faults.Fault;
import com.example.faults.Kind;

import java.nio.charset.StandardCharsets;
import java.time.Clock;
import java.time.Instant;
import java.time.ZoneOffset;
import java.util.Base64;
import java.util.Map;

import org.junit.jupiter.api.Test;

class AccessTokensTest {

  /** The Spec's two Tokens, Signed with spec-secret: one Lives until 2100, one Died at exp 1. */
  private static final String SPEC_LIVING = "eyJzdWIiOiJzdHVkZW50LXJlZ2lzdHJ5IiwiZXhwIjo0MTAyNDQ0ODAwfQ"
      + ".fgl9OtJeU8f90EIkPjELn34yWKjIbn3QcYn5UFoLAyY";
  private static final String SPEC_DEAD = "eyJzdWIiOiJzdHVkZW50LXJlZ2lzdHJ5IiwiZXhwIjoxfQ"
      + ".V8VQ8Ie_fo88z4EBJ0QbtIc5ONTOqQ_-l_WLuhwz5t4";

  private static final Instant ISSUED_AT = Instant.ofEpochSecond(1_000_000);

  private static AccessTokens clockedAt(String secret, Instant now) {
    return AccessTokens.buildAccessTokens(secret, 60, Clock.fixed(now, ZoneOffset.UTC));
  }

  private static String encodeClaims(String json) {
    return Base64.getUrlEncoder().withoutPadding().encodeToString(json.getBytes(StandardCharsets.UTF_8));
  }

  @Test
  void aTokenDiesOnlyAfterItsDeadline() {
    String token = clockedAt("secret", ISSUED_AT).issueAccessToken("student-registry");

    for (long after : new long[] {0, 59, 60}) {
      assertEquals("student-registry", clockedAt("secret", ISSUED_AT.plusSeconds(after)).readCaller(token));
    }
    assertEquals("student-registry",
        clockedAt("secret", ISSUED_AT.plusSeconds(60).plusMillis(999)).readCaller(token));

    assertThrows(Fault.class, () -> clockedAt("secret", ISSUED_AT.plusSeconds(61)).readCaller(token));
  }

  @Test
  void aMintedTokenCarriesCompactClaimsInTheSpecsOrder() {
    String token = clockedAt("secret", ISSUED_AT).issueAccessToken("student-registry");

    assertEquals(encodeClaims("{\"sub\":\"student-registry\",\"exp\":1000060}"), token.split("\\.")[0]);
  }

  @Test
  void theSpecsTokensReadTheSameInEveryRuntime() {
    AccessTokens spec = clockedAt("spec-secret", Instant.parse("2026-10-17T00:00:00Z"));

    assertEquals("student-registry", spec.readCaller(SPEC_LIVING));
    assertEquals(AccessTokens.TOKEN_IS_INVALID, assertThrows(Fault.class, () -> spec.readCaller(SPEC_DEAD)));
  }

  @Test
  void aTamperedClaimBreaksTheSignature() {
    String token = clockedAt("secret", ISSUED_AT).issueAccessToken("student-registry");
    String forged = encodeClaims("{\"sub\":\"admin\",\"exp\":99999999999}") + "." + token.split("\\.")[1];

    assertThrows(Fault.class, () -> clockedAt("secret", ISSUED_AT).readCaller(forged));
  }

  @Test
  void anotherSecretAndAnythingUnsignedAreRefused() {
    String token = clockedAt("secret", ISSUED_AT).issueAccessToken("student-registry");
    AccessTokens other = clockedAt("another", ISSUED_AT);

    Map.of("another secret", token, "unsigned", token.split("\\.")[0], "empty", "", "just a dot", ".",
        "not base64", "%%%.%%%", "three parts", token + ".x").forEach((name, candidate) ->
            assertEquals(AccessTokens.TOKEN_IS_INVALID, assertThrows(Fault.class, () -> other.readCaller(candidate)),
                name));
  }

  @Test
  void aTokenWithNoSubjectIsRefused() {
    AccessTokens tokens = clockedAt("spec-secret", ISSUED_AT);

    assertThrows(Fault.class, () -> tokens.readCaller(encodeClaims("{\"sub\":\"\",\"exp\":99999999999}") + ".x"));
  }

  @Test
  void anEmptySecretSignsNothingAndTrustsNothing() {
    assertThrows(Fault.class, () -> clockedAt("", ISSUED_AT).issueAccessToken("student-registry"));
    assertThrows(Fault.class, () -> clockedAt("secret", ISSUED_AT).issueAccessToken(""));
    assertThrows(Fault.class, () -> clockedAt("", ISSUED_AT).readCaller(SPEC_LIVING));
  }

  @Test
  void everyTokenFailureAnswersAsAnUnprovenCaller() {
    assertEquals(Kind.UNPROVEN_CALLER, AccessTokens.TOKEN_IS_INVALID.kind());
  }
}
