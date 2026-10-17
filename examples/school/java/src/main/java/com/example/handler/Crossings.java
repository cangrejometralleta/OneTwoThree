package com.example.handler;

import com.example.app.Validation;
import com.example.school.Course;
import com.example.school.CourseCode;
import com.example.school.CourseId;
import com.example.school.FullName;
import com.example.school.Rut;
import com.example.school.Student;
import com.example.school.StudentId;
import com.example.transport.Request;
import com.example.wire.CourseBody;
import com.example.wire.CourseView;
import com.example.wire.StudentBody;
import com.example.wire.StudentView;

import java.util.List;
import java.util.Map;
import java.util.Set;

/** The Wire never Becomes the Business by Accident: these are the Crossings. */
final class Crossings {

  private static final Set<String> STUDENT_FIELDS = Set.of("rut", "name", "age", "courseId");
  private static final Set<String> COURSE_FIELDS = Set.of("code", "name");

  private Crossings() {
  }

  static StudentBody readStudentBody(Request request) {
    Map<String, Object> sent = Validation.readJsonObject(request, STUDENT_FIELDS);

    return new StudentBody(Validation.readTextField(sent, "rut"), Validation.readTextField(sent, "name"),
        Validation.readWholeField(sent, "age"), Validation.readWholeField(sent, "courseId"));
  }

  static CourseBody readCourseBody(Request request) {
    Map<String, Object> sent = Validation.readJsonObject(request, COURSE_FIELDS);

    return new CourseBody(Validation.readTextField(sent, "code"), Validation.readTextField(sent, "name"));
  }

  static Student buildStudentRecord(StudentBody body, StudentId id) {
    return new Student(id, new Rut(body.rut()), new FullName(body.name()), body.age(), new CourseId(body.courseId()));
  }

  static StudentView renderStudentView(Student student) {
    return new StudentView(student.id().value(), student.rut().text(), student.name().text(), student.age(),
        student.course().value());
  }

  static List<StudentView> renderStudentViews(List<Student> students) {
    return students.stream().map(Crossings::renderStudentView).toList();
  }

  static Course buildCourseRecord(CourseBody body) {
    return new Course(CourseId.UNNUMBERED, new CourseCode(body.code()), new FullName(body.name()));
  }

  static CourseView renderCourseView(Course course) {
    return new CourseView(course.id().value(), course.code().text(), course.name().text());
  }

  static List<CourseView> renderCourseViews(List<Course> courses) {
    return courses.stream().map(Crossings::renderCourseView).toList();
  }
}
