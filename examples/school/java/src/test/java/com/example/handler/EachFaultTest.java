package com.example.handler;

import static org.junit.jupiter.api.Assertions.assertEquals;

import com.example.transport.Request;
import com.example.transport.Response;
import com.example.wire.FailureView;

import java.io.IOException;
import java.util.List;
import java.util.Map;
import java.util.function.Function;

import org.junit.jupiter.api.Test;

/**
 * Every Fault the Spec Names, Reached the Way a Caller Reaches it.
 * The Status and the Reason are Spelled here, never Read from the Fault under Test,
 * so a Declaration that Drifts Fails here instead of Agreeing with itself.
 */
class EachFaultTest {

  private static final String VALID = "{\"rut\":\"12345678-5\",\"name\":\"Ada\",\"age\":20,\"courseId\":1}";

  private record Case(String story, Function<TestingSchool, Response> arrive, int status, String reason) {
  }

  private static Function<TestingSchool, Response> enrolWith(String body) {
    return school -> school.callRoute("POST", "/students", Map.of(), body);
  }

  private static Function<TestingSchool, Response> showing(String pattern, String id) {
    return school -> school.callRoute("GET", pattern, Map.of("id", id), "");
  }

  @Test
  void eachFaultReachesTheEdgeWhole() {
    List<Case> cases = List.of(
        new Case("someone Enrols with no Name",
            enrolWith("{\"rut\":\"12345678-5\",\"name\":\"\",\"age\":20,\"courseId\":1}"), 400, "name is Empty"),
        new Case("someone Opens a Course with no Code",
            school -> school.callRoute("POST", "/courses", Map.of(), "{\"code\":\" \",\"name\":\"Physics\"}"),
            400, "code is Empty"),
        new Case("someone Enrols with a RUT that Fails its Digit",
            enrolWith("{\"rut\":\"12345678-9\",\"name\":\"Ada\",\"age\":20,\"courseId\":1}"), 400,
            "rut Fails its Check Digit"),
        new Case("someone Enrols below the Enrolment Age",
            enrolWith("{\"rut\":\"12345678-5\",\"name\":\"Ada\",\"age\":9,\"courseId\":1}"), 400,
            "age Must be 18 or more"),
        new Case("someone Asks for a Page that Counts backwards",
            school -> school.callRoute("GET", "/students", Map.of(), Map.of("page", "-1"), ""), 400, "page or size is not valid"),
        new Case("someone Sends a Body that is not JSON", enrolWith("{\"rut\":"), 400, "body is not valid JSON"),
        new Case("someone Sends a Body with an Identity of their own",
            enrolWith("{\"id\":5,\"rut\":\"12345678-5\",\"name\":\"Ada\",\"age\":20,\"courseId\":1}"), 400,
            "body is not valid JSON"),
        new Case("someone Asks for a Student by a Path that Holds no Number", showing("/students/{id}", "abc"), 400,
            "path Holds no valid Identity"),
        new Case("someone Arrives with no Token at all",
            school -> school.door.declareSchoolRoutes().routes().get(1).handle().handle(
                Request.EMPTY), 401, "token is Invalid or Expired"),
        new Case("someone Asks for a Student who never Enrolled", showing("/students/{id}", "42"), 404,
            "student not Found"),
        new Case("someone Asks for a Student beyond 64 bits", showing("/students/{id}", "99999999999999999999"), 404,
            "student not Found"),
        new Case("someone Enrols into a Course that never Opened",
            enrolWith("{\"rut\":\"12345678-5\",\"name\":\"Ada\",\"age\":20,\"courseId\":99}"), 404, "course not Found"),
        new Case("someone Knocks on a Door the Libretto does not Hold",
            school -> school.door.declareSchoolRoutes().unknown().handle(Request.EMPTY), 404,
            "route not Found"),
        new Case("someone Enrols with a RUT already Registered", school -> {
          enrolWith(VALID).apply(school);

          return enrolWith("{\"rut\":\"12345678-5\",\"name\":\"Grace\",\"age\":22,\"courseId\":1}").apply(school);
        }, 409, "rut is already Registered"),
        new Case("someone Enrols with a RUT the Registry Denies", school -> {
          school.front.registered = false;

          return enrolWith(VALID).apply(school);
        }, 422, "rut is not in the Registry"),
        new Case("someone Rewrites a Student nobody Holds with a RUT the Registry Denies", school -> {
          school.front.registered = false;

          return school.callRoute("PUT", "/students/{id}", Map.of("id", "7"), VALID);
        }, 422, "rut is not in the Registry"),
        new Case("someone Enrols while the Registry cannot Answer", school -> {
          school.front.down = new IOException("registry socket closed");

          return enrolWith(VALID).apply(school);
        }, 503, "registry cannot Answer now"),
        new Case("someone Drops a Student who already Left",
            school -> school.callRoute("DELETE", "/students/{id}", Map.of("id", "7"), ""), 404, "student not Found"),
        new Case("someone Rewrites an Enrolment that never Existed",
            school -> school.callRoute("PUT", "/students/{id}", Map.of("id", "7"), VALID), 404, "student not Found"));

    for (Case test : cases) {
      Response reply = test.arrive().apply(new TestingSchool());
      assertEquals(new Response(test.status(), new FailureView(test.reason())), reply, test.story());
    }
  }
}
