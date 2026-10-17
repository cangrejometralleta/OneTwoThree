package com.example.transport;

/** Handler Answers one Request. */
@FunctionalInterface
public interface Handler {
  Response handle(Request request);
}
