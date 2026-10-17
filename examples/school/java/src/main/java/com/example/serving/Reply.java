package com.example.serving;

import java.util.Optional;

/** Reply is what an Adapter Writes back: a Status, and a JSON Body unless there is none. */
public record Reply(int status, Optional<byte[]> json) {

  /** JSON_CONTENT_TYPE is the one Content-Type a Body Carries. */
  public static final String JSON_CONTENT_TYPE = "application/json";
}
