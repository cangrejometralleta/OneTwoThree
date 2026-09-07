package com.example.school;

import com.example.wire.CourseBody;
import com.example.wire.CourseView;
import com.example.wire.StudentBody;
import com.example.wire.StudentView;

import java.util.List;

/**
 * The Wire Package Holds Shapes; this File Holds the Crossing.
 * Promotion Lives with the Business Types it Promotes into.
 */
public final class Wire {

  private Wire() {
  }

  /** BuildStudentRecord Promotes untrusted Input into Business Types. */
  public static Student buildStudentRecord(StudentBody body, StudentId id) {
    return new Student(id, new Rut(text(body.rut())), new FullName(text(body.name())),
        body.age() == null ? 0 : body.age(), new CourseId(body.courseId() == null ? 0 : body.courseId()));
  }

  /** BuildCourseRecord Promotes untrusted Course Input. */
  public static Course buildCourseRecord(CourseBody body, CourseId id) {
    return new Course(id, new CourseCode(text(body.code())), new FullName(text(body.name())));
  }

  /** RenderStudentView Demotes a Business Student back to the Wire. */
  public static StudentView renderStudentView(Student student) {
    return new StudentView(student.id().value(), student.rut().value(), student.name().value(), student.age(),
        student.course().value());
  }

  /** RenderCourseView Demotes a Business Course back to the Wire. */
  public static CourseView renderCourseView(Course course) {
    return new CourseView(course.id().value(), course.code().value(), course.name().value());
  }

  /** RenderStudentViews Repeats the Move for a whole Page. */
  public static List<StudentView> renderStudentViews(List<Student> students) {
    return students.stream().map(Wire::renderStudentView).toList();
  }

  /** RenderCourseViews Repeats the Move for a whole Page. */
  public static List<CourseView> renderCourseViews(List<Course> courses) {
    return courses.stream().map(Wire::renderCourseView).toList();
  }

  private static String text(String value) {
    return value == null ? "" : value;
  }
}
