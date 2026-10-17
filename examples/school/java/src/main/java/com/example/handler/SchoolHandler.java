package com.example.handler;

import static java.net.HttpURLConnection.HTTP_CREATED;
import static java.net.HttpURLConnection.HTTP_NO_CONTENT;
import static java.net.HttpURLConnection.HTTP_OK;

import com.example.app.Answers;
import com.example.app.SchoolService;
import com.example.app.Telling;
import com.example.app.TokenIssuer;
import com.example.app.Validation;
import com.example.school.CourseId;
import com.example.school.StudentId;
import com.example.transport.Handler;
import com.example.transport.Libretto;
import com.example.transport.Request;
import com.example.transport.Route;
import com.example.wire.CourseView;
import com.example.wire.StudentView;
import com.example.wire.TokenView;

import java.util.List;

/**
 * SchoolHandler is the Door: it Orchestrates. It Decodes, Calls one Service method
 * and Answers with a Value or Fails. It Reaches no Provider and no Vendor.
 */
public record SchoolHandler(SchoolService school, TokenIssuer tokens) {

  /** CALLER_SUBJECT is the one Subject a Minted Token Names. */
  public static final String CALLER_SUBJECT = "student-registry";

  /** declareSchoolRoutes is the Libretto: every Door, its happy Status and its Guard. */
  public Libretto declareSchoolRoutes() {
    return new Libretto(List.of(
        new Route("POST", "/token", Answers.answerWith(HTTP_CREATED, this::mintAccessToken)),

        new Route("GET", "/students", guarded(HTTP_OK, this::listStudentRecords)),
        new Route("POST", "/students", guarded(HTTP_CREATED, this::addStudentRecord)),
        new Route("GET", "/students/{id}", guarded(HTTP_OK, this::showStudentRecord)),
        new Route("PUT", "/students/{id}", guarded(HTTP_OK, this::saveStudentRecord)),
        new Route("DELETE", "/students/{id}", guarded(HTTP_NO_CONTENT, this::dropStudentRecord)),

        new Route("GET", "/courses", guarded(HTTP_OK, this::listCourseRecords)),
        new Route("POST", "/courses", guarded(HTTP_CREATED, this::addCourseRecord)),
        new Route("GET", "/courses/{id}", guarded(HTTP_OK, this::showCourseRecord))),
        Answers.answerWith(HTTP_OK, this::refuseUnknownRoute));
  }

  /** guarded Names the Caller before the Script runs, then Crosses to a Handler. */
  public Handler guarded(int status, Telling tell) {
    return Answers.answerWith(status, Answers.requireProvenCaller(tokens, tell));
  }

  /** mintAccessToken Hands out a Token. */
  public TokenView mintAccessToken(Request request) {
    return new TokenView(tokens.issueAccessToken(CALLER_SUBJECT));
  }

  /** listStudentRecords Tells one Page of the Roll. */
  public List<StudentView> listStudentRecords(Request request) {
    return Crossings.renderStudentViews(school.listStudents(Validation.readPageRequest(request)));
  }

  /** addStudentRecord Enrols a new Student. */
  public StudentView addStudentRecord(Request request) {
    var body = Crossings.readStudentBody(request);

    return Crossings.renderStudentView(school.enrollStudent(Crossings.buildStudentRecord(body, StudentId.UNNUMBERED)));
  }

  /** showStudentRecord Tells the Story of one Student. */
  public StudentView showStudentRecord(Request request) {
    var id = new StudentId(Validation.readPathNumber(request));

    return Crossings.renderStudentView(school.readStudent(id));
  }

  /** saveStudentRecord Rewrites the Student the Path Names. */
  public StudentView saveStudentRecord(Request request) {
    var id = new StudentId(Validation.readPathNumber(request));
    var body = Crossings.readStudentBody(request);

    return Crossings.renderStudentView(school.saveStudent(Crossings.buildStudentRecord(body, id)));
  }

  /** dropStudentRecord Ends an Enrolment. It Answers no Body. */
  public Object dropStudentRecord(Request request) {
    school.deleteStudent(new StudentId(Validation.readPathNumber(request)));

    return null;
  }

  /** listCourseRecords Tells one Page of the Catalogue. */
  public List<CourseView> listCourseRecords(Request request) {
    return Crossings.renderCourseViews(school.listCourses(Validation.readPageRequest(request)));
  }

  /** addCourseRecord Opens a new Course. */
  public CourseView addCourseRecord(Request request) {
    return Crossings.renderCourseView(school.createCourse(Crossings.buildCourseRecord(Crossings.readCourseBody(request))));
  }

  /** showCourseRecord Tells the Story of one Course. */
  public CourseView showCourseRecord(Request request) {
    return Crossings.renderCourseView(school.readCourse(new CourseId(Validation.readPathNumber(request))));
  }

  /** refuseUnknownRoute Answers a Door the Libretto does not Hold, before any Guard. */
  public Object refuseUnknownRoute(Request request) {
    throw Validation.ROUTE_UNKNOWN;
  }
}
