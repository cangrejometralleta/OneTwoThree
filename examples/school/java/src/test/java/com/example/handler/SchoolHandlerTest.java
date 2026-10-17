package com.example.handler;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNull;

import com.example.transport.Request;
import com.example.transport.Response;
import com.example.transport.Route;
import com.example.wire.FailureView;
import com.example.wire.StudentView;

import java.io.IOException;
import java.util.List;
import java.util.Map;

import org.junit.jupiter.api.Test;

class SchoolHandlerTest {

  private static final String ENROL = "{\"rut\":\"12345678-5\",\"name\":\"Ada\",\"age\":20,\"courseId\":1}";
  private static final String REWRITE = "{\"rut\":\"12345678-5\",\"name\":\"Ada Lovelace\",\"age\":21,\"courseId\":1}";
  private static final Map<String, String> ONE = Map.of("id", "1");

  private final TestingSchool school = new TestingSchool();

  private record Story(String told, String method, String pattern, Map<String, String> path, String body,
      int status) {
  }

  @Test
  void eachRouteAnswersWithItsDeclaredStatus() {
    List<Story> stories = List.of(
        new Story("a Token is Minted", "POST", "/token", Map.of(), "", 201),
        new Story("someone Enrols", "POST", "/students", Map.of(), ENROL, 201),
        new Story("the Roll is Read", "GET", "/students", Map.of(), "", 200),
        new Story("one Student is Read", "GET", "/students/{id}", ONE, "", 200),
        new Story("an Enrolment is Rewritten", "PUT", "/students/{id}", ONE, REWRITE, 200),
        new Story("the Catalogue is Read", "GET", "/courses", Map.of(), "", 200),
        new Story("a Course Opens", "POST", "/courses", Map.of(), "{\"code\":\"FIS-201\",\"name\":\"Physics\"}", 201),
        new Story("one Course is Read", "GET", "/courses/{id}", ONE, "", 200),
        new Story("an Enrolment Ends", "DELETE", "/students/{id}", ONE, "", 204));

    for (Story story : stories) {
      Response reply = school.callRoute(story.method(), story.pattern(), story.path(), story.body());
      assertEquals(story.status(), reply.status(), story.told() + ": " + reply.body());
    }
  }

  @Test
  void theLibrettoDeclaresExactlyTheNineRoutesOfTheSpec() {
    assertEquals(9, school.door.declareSchoolRoutes().routes().size());
  }

  @Test
  void aRewriteChangesTheStudentAndADeleteLeavesNoBody() {
    school.callRoute("POST", "/students", Map.of(), ENROL);

    Response rewritten = school.callRoute("PUT", "/students/{id}", ONE, REWRITE);
    assertEquals(new StudentView(1, "12345678-5", "Ada Lovelace", 21, 1), rewritten.body());

    Response dropped = school.callRoute("DELETE", "/students/{id}", ONE, "");
    assertEquals(204, dropped.status());
    assertNull(dropped.body());
  }

  @Test
  void theViewEchoesTheRutInLowerCase() {
    school.callRoute("POST", "/students", Map.of(), "{\"rut\":\"10000013-K\",\"name\":\"Ada\",\"age\":20,\"courseId\":1}");

    assertEquals(new StudentView(1, "10000013-k", "Ada", 20, 1),
        school.callRoute("GET", "/students/{id}", ONE, "").body());
    assertEquals("10000013-k", school.front.asked.get(0).text());
  }

  @Test
  void anEmptyListAnswersAnEmptyArray() {
    assertEquals(List.of(), school.callRoute("GET", "/students", Map.of(), "").body());
  }

  @Test
  void aFailureReachesTheClientAsAFailureViewWithItsReason() {
    Response reply = school.callRoute("GET", "/students/{id}", Map.of("id", "42"), "");

    assertEquals(new Response(404, new FailureView("student not Found")), reply);
  }

  @Test
  void enrolmentSurvivesAFailedAnnouncement() {
    school.front.noticeFailure = new IOException("mail is down");

    assertEquals(201, school.callRoute("POST", "/students", Map.of(), ENROL).status());
  }

  @Test
  void everyRouteButTheTokenRefusesAnUnnamedCaller() {
    for (Route route : school.door.declareSchoolRoutes().routes()) {
      Response reply = route.handle().handle(new Request(ONE, Map.of(), Map.of(), new byte[0], ""));
      if (route.pattern().equals("/token")) {
        continue;
      }
      assertEquals(new Response(401, new FailureView("token is Invalid or Expired")), reply,
          route.method() + " " + route.pattern());
    }
  }

  @Test
  void anUnknownRouteAnswersBeforeTheGuard() {
    Response reply = school.door.declareSchoolRoutes().unknown().handle(Request.EMPTY);

    assertEquals(new Response(404, new FailureView("route not Found")), reply);
  }

  @Test
  void aGuardedHandlerReceivesTheNamedCaller() {
    String[] seen = new String[1];
    Response reply = school.door.guarded(200, request -> {
      seen[0] = request.caller();

      return "told";
    }).handle(new Request(Map.of(), Map.of(), Map.of("authorization",
        "Bearer " + school.tokens.issueAccessToken(SchoolHandler.CALLER_SUBJECT)), new byte[0], ""));

    assertEquals(200, reply.status());
    assertEquals("student-registry", seen[0]);
  }

  @Test
  void anUncontrolledFailureAnswersFiveHundredAndSaysNothingOfTheDriver() {
    Response reply = school.door.guarded(200, request -> {
      throw new IllegalStateException("the driver Broke at /var/db");
    }).handle(new Request(Map.of(), Map.of(), Map.of("authorization",
        "Bearer " + school.tokens.issueAccessToken(SchoolHandler.CALLER_SUBJECT)), new byte[0], ""));

    assertEquals(new Response(500, new FailureView("internal error")), reply);
  }
}
