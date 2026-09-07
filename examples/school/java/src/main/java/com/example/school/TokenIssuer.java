package com.example.school;

/**
 * TokenIssuer Mints and Reads a Bearer Token.
 * Fulfilled by tokens.AccessTokens, so no third Party Enters for this.
 */
public interface TokenIssuer {

  String issueAccessToken(String subject);

  String readAccessToken(String token);
}
