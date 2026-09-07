package com.example.app;

import com.example.school.TokenIssuer;

/** The Caller is Named before the Story Starts. */
public final class Context {

  private Context() {
  }

  /**
   * RequireProvenCaller Names the Caller before the Story Starts.
   * A Handler behind this Reads request.caller() and Trusts it.
   */
  public static Telling requireProvenCaller(TokenIssuer tokens, Telling tell) {
    return request -> tell.tell(request.namedFor(tokens.readAccessToken(request.token())));
  }
}
