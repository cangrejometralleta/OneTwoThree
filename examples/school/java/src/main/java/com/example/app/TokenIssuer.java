package com.example.app;

/** TokenIssuer Mints a Bearer Token and Reads one for the Guard. */
public interface TokenIssuer extends CallerReader {
  String issueAccessToken(String subject);
}
