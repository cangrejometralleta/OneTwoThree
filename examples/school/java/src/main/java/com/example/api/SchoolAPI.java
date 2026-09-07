package com.example.api;

import com.example.app.Answers;
import com.example.app.Context;
import com.example.app.Telling;
import com.example.app.Validation;
import com.example.school.Course;
import com.example.school.Page;
import com.example.school.CourseStore;
import com.example.school.StudentStore;
import com.example.school.TokenIssuer;
import com.example.school.Student;
import com.example.school.CourseId;
import com.example.school.StudentId;
import com.example.school.Wire;
import com.example.transport.Handler;
import com.example.transport.Request;
import com.example.transport.Route;
import com.example.wire.CourseBody;
import com.example.wire.StudentBody;

import static java.net.HttpURLConnection.HTTP_CREATED;
import static java.net.HttpURLConnection.HTTP_NO_CONTENT;
import static java.net.HttpURLConnection.HTTP_OK;

import java.util.List;
import java.util.Map;

/**
 * SchoolAPI Holds the Cast every Story Needs.
 * Three Providers, no Libraries: a Test can Hand it three Fakes.
 */
public record SchoolAPI(StudentStore students, CourseStore courses, TokenIssuer tokens) {

  /**
   * DeclareSchoolRoutes is the Libretto.
   * Read it once and you Know the whole Service.
   */
  public List<Route> declareSchoolRoutes() {
    return List.of(
        new Route("POST", "/token", Answers.answerWith(HTTP_CREATED, this::mintAccessToken)),

        new Route("GET", "/students", guarded(HTTP_OK, this::listStudentRecords)),
        new Route("POST", "/students", guarded(HTTP_CREATED, this::addStudentRecord)),
        new Route("GET", "/students/{id}", guarded(HTTP_OK, this::showStudentRecord)),
        new Route("PUT", "/students/{id}", guarded(HTTP_OK, this::saveStudentRecord)),
        new Route("DELETE", "/students/{id}", guarded(HTTP_NO_CONTENT, this::dropStudentRecord)),

        new Route("GET", "/courses", guarded(HTTP_OK, this::listCourseRecords)),
        new Route("POST", "/courses", guarded(HTTP_CREATED, this::addCourseRecord)),
        new Route("GET", "/courses/{id}", guarded(HTTP_OK, this::showCourseRecord)));
  }

  /**
   * Guarded Says the same two Things about a Route every time:
   * Name the Caller first, then Answer with this Status.
   */
  public Handler guarded(int status, Telling tell) {
    return Answers.answerWith(status, Context.requireProvenCaller(tokens, tell));
  }

  /** MintAccessToken Hands out a Token. */
  public Object mintAccessToken(Request request) {
    return Map.of("token", tokens.issueAccessToken("student-registry"));
  }

  /** ListStudentRecords Tells one Page of the Roll. */
  public Object listStudentRecords(Request request) {
    Page page = Validation.readPageRequest(request);

    return Wire.renderStudentViews(students.selectStudentPage(page));
  }

  /** ShowStudentRecord Tells the Story of one Student. */
  public Object showStudentRecord(Request request) {
    StudentId id = new StudentId(Validation.readPathNumber(request));

    return Wire.renderStudentView(students.selectStudentRow(id));
  }

  /** AddStudentRecord Enrols someone new. */
  public Object addStudentRecord(Request request) {
    Student student = readStudentBody(request, new StudentId(0));

    return Wire.renderStudentView(students.insertStudentRow(student));
  }

  /** SaveStudentRecord Rewrites an Enrolment that already Exists. */
  public Object saveStudentRecord(Request request) {
    StudentId id = new StudentId(Validation.readPathNumber(request));
    Student student = readStudentBody(request, id);

    return Wire.renderStudentView(students.updateStudentRow(student));
  }

  /** DropStudentRecord Ends an Enrolment. */
  public Object dropStudentRecord(Request request) {
    students.deleteStudentRow(new StudentId(Validation.readPathNumber(request)));

    return null;
  }

  /** ListCourseRecords Tells one Page of the Catalogue. */
  public Object listCourseRecords(Request request) {
    Page page = Validation.readPageRequest(request);

    return Wire.renderCourseViews(courses.selectCoursePage(page));
  }

  /** ShowCourseRecord Tells the Story of one Course. */
  public Object showCourseRecord(Request request) {
    CourseId id = new CourseId(Validation.readPathNumber(request));

    return Wire.renderCourseView(courses.selectCourseRow(id));
  }

  /** AddCourseRecord Opens a new Course. */
  public Object addCourseRecord(Request request) {
    CourseBody body = Validation.readJsonBody(request, CourseBody.class);
    Course course = Wire.buildCourseRecord(body, new CourseId(0));

    return Wire.renderCourseView(courses.insertCourseRow(course));
  }

  /** ReadStudentBody Decodes the Wire, Validates it, then Checks the Course. */
  private Student readStudentBody(Request request, StudentId id) {
    StudentBody body = Validation.readJsonBody(request, StudentBody.class);
    Student student = Wire.buildStudentRecord(body, id);

    student.checkStudentRecord();
    courses.selectCourseRow(student.course());

    return student;
  }
}
