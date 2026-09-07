package com.example.transport;

/** Handler is the only Signature the Business Layer Knows. */
@FunctionalInterface
public interface Handler {
  Response handle(Request request);
}
