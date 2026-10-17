package com.example.app;

import com.example.transport.Request;

/** Telling is a Script: it Answers with a Value or Fails, and builds no Reply. */
@FunctionalInterface
public interface Telling {
  Object tell(Request request) throws Exception;
}
