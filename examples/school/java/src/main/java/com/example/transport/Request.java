package com.example.transport;

import java.util.Map;

/**
 * Request is what a Handler Receives.
 * No Framework Type Appears here.
 *
 * <p>Caller is Empty until the Application Names it.
 */
public record Request(Map<String, String> path, Map<String, String> query, String token, byte[] body, String caller) {

  public static Request empty() {
    return new Request(Map.of(), Map.of(), "", new byte[0], "");
  }

  public Request namedFor(String caller) {
    return new Request(path, query, token, body, caller);
  }

  public String pathValue(String name) {
    return path.getOrDefault(name, "");
  }

  public String queryValue(String name) {
    return query.getOrDefault(name, "");
  }
}
