package com.example.tokens;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

import static java.net.HttpURLConnection.HTTP_UNAUTHORIZED;

import com.example.faults.Fault;

import org.junit.jupiter.api.Test;

import java.time.Instant;
import java.util.concurrent.atomic.AtomicReference;

class AccessTokensTest {

  /** The Comparison that Broke in the Spring Version, Pinned here. */
  @Test
  void accessTokenDiesOnlyAfterItsDeadline() {
    AtomicReference<Instant> clock = new AtomicReference<>(Instant.now());
    AccessTokens tokens = AccessTokens.buildAccessTokensReading("s", 60, clock::get);

    String token = tokens.issueAccessToken("registry");
    assertEquals("registry", tokens.readAccessToken(token));

    clock.set(clock.get().plusSeconds(120));
    assertThrows(Fault.class, () -> tokens.readAccessToken(token));
  }

  /** Every Way a Token Fails Arrives as the same Answer. */
  @Test
  void everyTokenFailureAnswersUnauthorised() {
    AccessTokens tokens = AccessTokens.buildAccessTokens("s", 60);
    String signed = tokens.issueAccessToken("registry");

    for (String broken : new String[] {"", "hello", "registry.99999999999", "admin" + signed.substring(8)}) {
      Fault fault = assertThrows(Fault.class, () -> tokens.readAccessToken(broken));
      assertEquals(HTTP_UNAUTHORIZED, fault.status());
    }
  }

  @Test
  void issueAccessTokenRefusesAnEmptySecret() {
    assertThrows(Fault.class, () -> AccessTokens.buildAccessTokens("", 60).issueAccessToken("registry"));
  }
}
