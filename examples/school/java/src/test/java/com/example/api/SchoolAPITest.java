package com.example.api;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.fail;

import static java.net.HttpURLConnection.HTTP_BAD_REQUEST;
import static java.net.HttpURLConnection.HTTP_CONFLICT;
import static java.net.HttpURLConnection.HTTP_CREATED;
import static java.net.HttpURLConnection.HTTP_INTERNAL_ERROR;
import static java.net.HttpURLConnection.HTTP_NOT_FOUND;
import static java.net.HttpURLConnection.HTTP_NO_CONTENT;
import static java.net.HttpURLConnection.HTTP_OK;
import static java.net.HttpURLConnection.HTTP_UNAUTHORIZED;

import com.example.faults.Fault;
import com.example.school.Course;
import com.example.school.CourseCode;
import com.example.school.CourseId;
import com.example.school.CourseStore;
import com.example.school.Errors;
import com.example.school.FullName;
import com.example.school.Page;
import com.example.school.Student;
import com.example.school.StudentId;
import com.example.school.StudentStore;
import com.example.tokens.AccessTokens;
import com.example.transport.Request;
import com.example.transport.Response;
import com.example.transport.Route;

import org.junit.jupiter.api.Test;

import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

class SchoolAPITest {

  /**
   * FakeSchool Stands in for JPA.
   * No Database, no Framework, no Port: the Providers Allow it.
   */
  static class FakeSchool implements StudentStore, CourseStore {

    private final Map<Long, Student> students = new LinkedHashMap<>();
    private final Map<Long, Course> courses = new LinkedHashMap<>();

    FakeSchool() {
      courses.put(1L, new Course(new CourseId(1), new CourseCode("MAT-101"), new FullName("Algebra")));
    }

    /** InsertStudentRow Keeps the RUT unique, the Way a real Index would. */
    @Override
    public Student insertStudentRow(Student student) {
      for (Student enrolled : students.values()) {
        if (enrolled.rut().equals(student.rut())) {
          throw Errors.RUT_TAKEN;
        }
      }

      Student numbered = new Student(new StudentId(students.size() + 1), student.rut(), student.name(), student.age(),
          student.course());
      students.put(numbered.id().value(), numbered);

      return numbered;
    }

    @Override
    public Student selectStudentRow(StudentId id) {
      Student student = students.get(id.value());
      if (student == null) {
        throw Errors.STUDENT_UNKNOWN;
      }

      return student;
    }

    @Override
    public Student updateStudentRow(Student student) {
      selectStudentRow(student.id());
      students.put(student.id().value(), student);

      return student;
    }

    @Override
    public void deleteStudentRow(StudentId id) {
      selectStudentRow(id);
      students.remove(id.value());
    }

    @Override
    public List<Student> selectStudentPage(Page page) {
      return new ArrayList<>(students.values());
    }

    @Override
    public Course insertCourseRow(Course course) {
      Course numbered = new Course(new CourseId(courses.size() + 1), course.code(), course.name());
      courses.put(numbered.id().value(), numbered);

      return numbered;
    }

    @Override
    public Course selectCourseRow(CourseId id) {
      Course course = courses.get(id.value());
      if (course == null) {
        throw Errors.COURSE_UNKNOWN;
      }

      return course;
    }

    @Override
    public List<Course> selectCoursePage(Page page) {
      return new ArrayList<>(courses.values());
    }
  }

  /** BuildTestingSchool Hands the API three Fakes. */
  static SchoolAPI buildTestingSchool() {
    FakeSchool fake = new FakeSchool();

    return new SchoolAPI(fake, fake, AccessTokens.buildAccessTokens("test", 60));
  }

  static Request requestWith(String body, Map<String, String> path, Map<String, String> query) {
    return new Request(path, query, "", body.getBytes(StandardCharsets.UTF_8), "");
  }

  static final String GOOD = "{\"rut\":\"12345678-5\",\"name\":\"Ada\",\"age\":20,\"courseId\":1}";

  /**
   * Every Fault the Service Declares, Reached the Way a Caller Reaches it.
   * Each Name is a Use Case, because a controlled Failure is one.
   */
  @Test
  void eachFaultReachesTheEdgeWhole() {
    checkFault("someone Enrols with no Name",
        api -> api.addStudentRecord(requestWith("{\"rut\":\"12345678-5\",\"name\":\"\",\"age\":20,\"courseId\":1}",
            Map.of(), Map.of())),
        Errors.NAME_IS_EMPTY, HTTP_BAD_REQUEST);

    checkFault("someone Enrols with a RUT that Fails its Digit",
        api -> api.addStudentRecord(requestWith("{\"rut\":\"12345678-9\",\"name\":\"Ada\",\"age\":20,\"courseId\":1}",
            Map.of(), Map.of())),
        Errors.RUT_IS_INVALID, HTTP_BAD_REQUEST);

    checkFault("someone Enrols below the Enrolment Age",
        api -> api.addStudentRecord(requestWith("{\"rut\":\"12345678-5\",\"name\":\"Ada\",\"age\":9,\"courseId\":1}",
            Map.of(), Map.of())),
        Errors.AGE_IS_TOO_LOW, HTTP_BAD_REQUEST);

    checkFault("someone Enrols into a Course that never Opened",
        api -> api.addStudentRecord(requestWith("{\"rut\":\"12345678-5\",\"name\":\"Ada\",\"age\":20,\"courseId\":99}",
            Map.of(), Map.of())),
        Errors.COURSE_UNKNOWN, HTTP_NOT_FOUND);

    checkFault("someone Sends a Body that is not JSON",
        api -> api.addStudentRecord(requestWith("{\"rut\":", Map.of(), Map.of())),
        com.example.app.Validation.BODY_IS_BROKEN, HTTP_BAD_REQUEST);

    checkFault("someone Enrols with a RUT already Registered", api -> {
      api.addStudentRecord(requestWith(GOOD, Map.of(), Map.of()));

      return api.addStudentRecord(requestWith(GOOD, Map.of(), Map.of()));
    }, Errors.RUT_TAKEN, HTTP_CONFLICT);

    checkFault("someone Asks for a Student by a Path that Holds no Number",
        api -> api.showStudentRecord(requestWith("", Map.of("id", "abc"), Map.of())),
        com.example.app.Validation.PATH_IS_BROKEN, HTTP_BAD_REQUEST);

    checkFault("someone Asks for a Student who never Enrolled",
        api -> api.showStudentRecord(requestWith("", Map.of("id", "42"), Map.of())),
        Errors.STUDENT_UNKNOWN, HTTP_NOT_FOUND);

    checkFault("someone Asks for a Course that never Opened",
        api -> api.showCourseRecord(requestWith("", Map.of("id", "99"), Map.of())),
        Errors.COURSE_UNKNOWN, HTTP_NOT_FOUND);

    checkFault("someone Asks for a Page that Counts backwards",
        api -> api.listStudentRecords(requestWith("", Map.of(), Map.of("page", "-1"))),
        Errors.PAGE_IS_INVALID, HTTP_BAD_REQUEST);

    checkFault("someone Asks for a Page Numbered with a Word",
        api -> api.listStudentRecords(requestWith("", Map.of(), Map.of("page", "abc"))),
        Errors.PAGE_IS_INVALID, HTTP_BAD_REQUEST);

    checkFault("someone Asks for a Size Measured in Words",
        api -> api.listCourseRecords(requestWith("", Map.of(), Map.of("size", "many"))),
        Errors.PAGE_IS_INVALID, HTTP_BAD_REQUEST);

    checkFault("someone Drops a Student who already Left",
        api -> api.dropStudentRecord(requestWith("", Map.of("id", "7"), Map.of())),
        Errors.STUDENT_UNKNOWN, HTTP_NOT_FOUND);

    checkFault("someone Arrives with no Token at all",
        api -> api.guarded(HTTP_OK, api::listStudentRecords).handle(Request.empty()),
        AccessTokens.TOKEN_IS_INVALID, HTTP_UNAUTHORIZED);
  }

  /**
   * CheckFault Reads both Halves: the Fault that Arrived and the Answer it Carries.
   * The Number is Spelled here, never Read from the Fault under Test.
   */
  private void checkFault(String story, java.util.function.Function<SchoolAPI, Object> arrive, Fault want, int status) {
    SchoolAPI api = buildTestingSchool();
    try {
      Object answered = arrive.apply(api);
      if (answered instanceof Response reply && reply.status() == status) {
        return;
      }
      fail(story + ": wanted " + want.getMessage() + ", got " + answered);
    } catch (Fault got) {
      assertEquals(want.getMessage(), got.getMessage(), story);
      assertEquals(status, got.status(), story);
    }
  }

  /** The happy Status Lives in the Route, so a Test Reads it from there. */
  @Test
  void eachRouteAnswersWithItsDeclaredStatus() {
    SchoolAPI api = buildTestingSchool();
    String token = api.tokens().issueAccessToken("student-registry");

    assertEquals(HTTP_CREATED, callRoute(api, "POST", "/token", "", Map.of(), "").status());
    assertEquals(HTTP_CREATED, callRoute(api, "POST", "/students", GOOD, Map.of(), token).status());
    assertEquals(HTTP_OK, callRoute(api, "GET", "/students", "", Map.of(), token).status());
    assertEquals(HTTP_OK, callRoute(api, "GET", "/students/{id}", "", Map.of("id", "1"), token).status());
    assertEquals(HTTP_NO_CONTENT, callRoute(api, "DELETE", "/students/{id}", "", Map.of("id", "1"), token).status());
  }

  /** Every Route but the Token Names its Caller first. */
  @Test
  void everyRouteButTheTokenRefusesAnUnnamedCaller() {
    SchoolAPI api = buildTestingSchool();

    for (Route route : api.declareSchoolRoutes()) {
      int wanted = route.pattern().equals("/token") ? HTTP_CREATED : HTTP_UNAUTHORIZED;
      Response reply = route.handle().handle(requestWith("", Map.of("id", "1"), Map.of()));

      assertEquals(wanted, reply.status(), route.method() + " " + route.pattern());
    }
  }

  /** A Guarded Handler Reads a Caller it never had to Check. */
  @Test
  void guardedHandlerReceivesTheNamedCaller() {
    SchoolAPI api = buildTestingSchool();
    String token = api.tokens().issueAccessToken("student-registry");
    StringBuilder seen = new StringBuilder();

    api.guarded(HTTP_OK, request -> {
      seen.append(request.caller());

      return null;
    }).handle(new Request(Map.of(), Map.of(), token, new byte[0], ""));

    assertEquals("student-registry", seen.toString());
  }

  /** An uncontrolled Store Failure never Becomes the Caller's Fault. */
  @Test
  void uncontrolledFailureAnswersFiveHundred() {
    FakeSchool broken = new FakeSchool() {
      @Override
      public List<Student> selectStudentPage(Page page) {
        throw new IllegalStateException("the driver Closed the Connection");
      }
    };
    SchoolAPI api = new SchoolAPI(broken, broken, AccessTokens.buildAccessTokens("test", 60));
    String token = api.tokens().issueAccessToken("student-registry");

    assertEquals(HTTP_INTERNAL_ERROR, callRoute(api, "GET", "/students", "", Map.of(), token).status());
  }

  private Response callRoute(SchoolAPI api, String method, String pattern, String body, Map<String, String> path,
      String token) {
    for (Route route : api.declareSchoolRoutes()) {
      if (route.method().equals(method) && route.pattern().equals(pattern)) {
        return route.handle().handle(new Request(path, Map.of(), token, body.getBytes(StandardCharsets.UTF_8), ""));
      }
    }

    throw new IllegalStateException("no Route Answers " + method + " " + pattern);
  }
}
