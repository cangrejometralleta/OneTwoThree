package com.example.handler;

import com.example.app.FrontDesk;
import com.example.app.SchoolService;
import com.example.app.Shelf;
import com.example.school.SchoolRules;
import com.example.tokens.AccessTokens;
import com.example.transport.Request;
import com.example.transport.Response;
import com.example.transport.Route;

import java.nio.charset.StandardCharsets;
import java.time.Clock;
import java.util.Map;

/** TestingSchool Hands the Door its Fakes: no Database, no Framework, no Port. */
final class TestingSchool {

  final Shelf books = new Shelf();
  final FrontDesk front = new FrontDesk();
  final AccessTokens tokens = AccessTokens.buildAccessTokens("test", 60, Clock.systemUTC());
  final SchoolHandler door = new SchoolHandler(
      new SchoolService(new SchoolRules(18), books, books, front, front, System.getLogger("school.test")), tokens);

  /** callRoute Arrives the Way a Client Arrives: through the Libretto, the Answer and the Guard. */
  Response callRoute(String method, String pattern, Map<String, String> path, String body) {
    return callRoute(method, pattern, path, Map.of(), body);
  }

  /** callRoute Carries a Query too, for the list Routes. */
  Response callRoute(String method, String pattern, Map<String, String> path, Map<String, String> query,
      String body) {
    Request request = new Request(path, query, Map.of("authorization", "Bearer " + tokens.issueAccessToken(
        SchoolHandler.CALLER_SUBJECT)), body.getBytes(StandardCharsets.UTF_8), "");

    for (Route route : door.declareSchoolRoutes().routes()) {
      if (route.method().equals(method) && route.pattern().equals(pattern)) {
        return route.handle().handle(request);
      }
    }

    throw new AssertionError("no Route Answers " + method + " " + pattern);
  }
}
