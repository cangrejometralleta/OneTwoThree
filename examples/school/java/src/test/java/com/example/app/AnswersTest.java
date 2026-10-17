package com.example.app;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;

import com.example.faults.Fault;
import com.example.transport.Request;
import com.example.transport.Response;
import com.example.wire.FailureView;

import java.util.Map;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicReference;

import org.junit.jupiter.api.Test;

class AnswersTest {

  private static final CallerReader FIXED_READER = token -> {
    if (!token.equals("good")) {
      throw Fault.refuseUnprovenCaller("token is Invalid or Expired");
    }

    return "student-registry";
  };

  private static Request bearing(String header) {
    return new Request(Map.of(), Map.of(), Map.of("authorization", header), new byte[0], "");
  }

  @Test
  void theCallerIsNamedBeforeTheStoryStartsWhateverTheCaseOfTheScheme() throws Exception {
    AtomicReference<String> seen = new AtomicReference<>();
    Telling guarded = Answers.requireProvenCaller(FIXED_READER, request -> {
      seen.set(request.caller());

      return "told";
    });

    for (String header : new String[] {"Bearer good", "bearer good", "BEARER good"}) {
      seen.set(null);
      guarded.tell(bearing(header));
      assertEquals("student-registry", seen.get(), header);
    }
  }

  @Test
  void theStoryNeverStartsForAStranger() {
    AtomicBoolean started = new AtomicBoolean();
    Telling guarded = Answers.requireProvenCaller(FIXED_READER, request -> {
      started.set(true);

      return null;
    });

    Response reply = Answers.answerWith(200, guarded).handle(bearing("Bearer forged"));

    assertEquals(401, reply.status());
    assertFalse(started.get());
  }

  @Test
  void anAnswerCarriesTheDeclaredStatusOrTheFaultsOwn() {
    assertEquals(new Response(201, "made"), Answers.answerWith(201, request -> "made").handle(Request.EMPTY));

    Response failing = Answers.answerWith(200, request -> {
      throw Fault.reportMissingRecord("gone");
    }).handle(Request.EMPTY);
    assertEquals(new Response(404, new FailureView("gone")), failing);

    Response ours = Answers.answerWith(200, request -> {
      throw new IllegalStateException("driver at /var/db broke");
    }).handle(Request.EMPTY);
    assertEquals(new Response(500, new FailureView("internal error")), ours);
  }
}
