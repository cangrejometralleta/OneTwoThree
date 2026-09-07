package com.example.app;

import com.example.faults.Fault;
import com.example.transport.Handler;
import com.example.transport.Response;

import java.util.Map;

/** The Crossing: the only Place in the Program that Builds a Reply. */
public final class Answers {

  private Answers() {
  }

  /**
   * AnswerWith Turns a Telling into a Handler the Transport can Mount.
   * The Route Declares the happy Status; a Fault Declares its own.
   */
  public static Handler answerWith(int status, Telling tell) {
    return request -> {
      try {
        return new Response(status, tell.tell(request));
      } catch (RuntimeException failure) {
        return new Response(Fault.readFaultStatus(failure), Map.of("error", reasonOf(failure)));
      }
    };
  }

  private static String reasonOf(RuntimeException failure) {
    return failure.getMessage() == null ? "the Service Failed" : failure.getMessage();
  }
}
