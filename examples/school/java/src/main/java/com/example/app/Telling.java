package com.example.app;

import com.example.transport.Request;

/**
 * A Telling Speaks Business only: it Answers with a Value, or it Fails.
 * It Names no Status and Builds no Reply.
 */
@FunctionalInterface
public interface Telling {
  Object tell(Request request);
}
