package com.example.app;

import com.example.faults.Fault;
import com.example.transport.Handler;
import com.example.transport.Response;
import com.example.wire.FailureView;

/** Answers Holds the two Crossings every Route Shares: the Answer and the Guard. */
public final class Answers {

  private Answers() {
  }

  /**
   * answerWith is the one Crossing from a Telling to a Handler.
   * The Route Declares the happy Status; a Fault Declares its own.
   */
  public static Handler answerWith(int status, Telling tell) {
    return request -> {
      try {
        return new Response(status, tell.tell(request));
      } catch (Exception failure) {
        return new Response(Statuses.readFaultStatus(failure), new FailureView(Fault.readFaultReason(failure)));
      }
    };
  }

  /**
   * requireProvenCaller Names the Caller before the Story starts.
   * A Script behind it Reads the Caller and Trusts it.
   */
  public static Telling requireProvenCaller(CallerReader reader, Telling tell) {
    return request -> tell.tell(request.withCaller(reader.readCaller(Validation.readBearerToken(request))));
  }
}
