package com.example.app;

/** CallerReader Names the Caller a Token Proves, or Fails. */
public interface CallerReader {
  String readCaller(String token);
}
