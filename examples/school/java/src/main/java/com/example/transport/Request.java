package com.example.transport;

import java.util.Map;

/**
 * Request is what the Door Hands the Script: no Framework Type Crosses it.
 * Header Names Arrive in lower case; the Body is at most one Byte past the Bound.
 */
public record Request(Map<String, String> path, Map<String, String> query, Map<String, String> headers,
    byte[] body, String caller) {

  /** MAX_BODY_BYTES Bounds what a Client may Send: one mebibyte. */
  public static final int MAX_BODY_BYTES = 1 << 20;

  public Request {
    path = Map.copyOf(path);
    query = Map.copyOf(query);
    headers = Map.copyOf(headers);
    body = body.clone();
  }

  /** EMPTY Holds nothing: no Path, no Query, no Header and no Body. */
  public static final Request EMPTY = new Request(Map.of(), Map.of(), Map.of(), new byte[0], "");

  /** withCaller Names the Caller the Guard Proved. */
  public Request withCaller(String named) {
    return new Request(path, query, headers, body, named);
  }

  /** withPath Hands the Route's Names to the Script. */
  public Request withPath(Map<String, String> names) {
    return new Request(names, query, headers, body, caller);
  }

  @Override
  public byte[] body() {
    return body.clone();
  }
}
