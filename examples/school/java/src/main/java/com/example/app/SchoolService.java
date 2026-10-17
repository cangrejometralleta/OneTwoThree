package com.example.app;

import com.example.school.Course;
import com.example.school.CourseId;
import com.example.school.Errors;
import com.example.school.Page;
import com.example.school.Rut;
import com.example.school.SchoolRules;
import com.example.school.Student;
import com.example.school.StudentId;

import java.lang.System.Logger;
import java.lang.System.Logger.Level;
import java.util.List;

/**
 * SchoolService Joins the Providers behind each School Use Case.
 * It never Learns that two Questions may have one Answerer.
 */
public record SchoolService(SchoolRules rules, StudentStore students, CourseStore courses,
    RutRegistry registry, EnrollmentNotifier notices, Logger logger) {

  /** listStudents Reads one Page of the Roll. */
  public List<Student> listStudents(Page page) {
    return students.selectStudentPage(page);
  }

  /** readStudent Finds one Enrolment, or it is Unknown. */
  public Student readStudent(StudentId id) {
    return students.selectStudentRow(id).orElseThrow(() -> Errors.STUDENT_UNKNOWN);
  }

  /**
   * enrollStudent Validates and Keeps one Enrolment, then Announces it.
   * An Announcement that Fails never Fails the Enrolment: the Student Exists.
   */
  public Student enrollStudent(Student student) {
    checkStudent(student);
    Student kept = students.insertStudentRow(student);
    announceEnrollment(kept);

    return kept;
  }

  /** saveStudent Validates and Rewrites one Enrolment. It Announces nothing. */
  public Student saveStudent(Student student) {
    checkStudent(student);
    if (!students.updateStudentRow(student)) {
      throw Errors.STUDENT_UNKNOWN;
    }

    return student;
  }

  /** deleteStudent Ends one Enrolment, once. */
  public void deleteStudent(StudentId id) {
    if (!students.deleteStudentRow(id)) {
      throw Errors.STUDENT_UNKNOWN;
    }
  }

  /** listCourses Reads one Page of the Catalogue. */
  public List<Course> listCourses(Page page) {
    return courses.selectCoursePage(page);
  }

  /** readCourse Finds one Course, or it is Unknown. */
  public Course readCourse(CourseId id) {
    return courses.selectCourseRow(id).orElseThrow(() -> Errors.COURSE_UNKNOWN);
  }

  /** createCourse Validates and Keeps one Course. */
  public Course createCourse(Course course) {
    course.checkCourseRecord();

    return courses.insertCourseRow(course);
  }

  /** checkStudent Asks in the Order the Spec Names: Form, Course, Registry. */
  private void checkStudent(Student student) {
    student.checkStudentRecord(rules);
    readCourse(student.course());
    confirmRut(student.rut());
  }

  /** confirmRut Asks the Registry. An Unknown is never a Yes. */
  private void confirmRut(Rut rut) {
    if (registry == null) {
      throw Errors.REGISTRY_UNAVAILABLE;
    }
    if (!askRegistry(rut)) {
      throw Errors.RUT_UNREGISTERED;
    }
  }

  /** askRegistry Turns a Registry that cannot Answer into Registry Unavailable. */
  private boolean askRegistry(Rut rut) {
    try {
      return registry.confirmRut(rut);
    } catch (Exception silent) {
      throw Errors.REGISTRY_UNAVAILABLE;
    }
  }

  /**
   * announceEnrollment Keeps a Notice Failure inside the Answer: the Caller
   * cannot Act on a Notice that did not Leave, so it is Logged and no more.
   */
  private void announceEnrollment(Student student) {
    if (notices == null) {
      return;
    }

    try {
      notices.announceEnrollment(student);
    } catch (Exception failed) {
      logger.log(Level.WARNING, "enrollment announcement failed for student " + student.id().value());
    }
  }
}
