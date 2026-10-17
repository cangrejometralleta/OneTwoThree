package com.example.app;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

import com.example.faults.Fault;
import com.example.school.Course;
import com.example.school.CourseCode;
import com.example.school.CourseId;
import com.example.school.Errors;
import com.example.school.FullName;
import com.example.school.Page;
import com.example.school.Rut;
import com.example.school.SchoolRules;
import com.example.school.Student;
import com.example.school.StudentId;

import java.io.IOException;

import org.junit.jupiter.api.Test;

class SchoolServiceTest {

  private final Shelf books = new Shelf();
  private final FrontDesk front = new FrontDesk();
  private final SchoolService service = buildAdmission(front);

  private SchoolService buildAdmission(FrontDesk desk) {
    return new SchoolService(new SchoolRules(18), books, books, desk, desk, System.getLogger("school.test"));
  }

  private static Student readyStudent() {
    return new Student(StudentId.UNNUMBERED, new Rut("11111111-1"), new FullName("Ana"), 20, new CourseId(1));
  }

  private static Student withCourse(Student student, long course) {
    return new Student(student.id(), student.rut(), student.name(), student.age(), new CourseId(course));
  }

  @Test
  void enrollStudentAnnouncesOnceTheStoreKeeps() {
    Student kept = service.enrollStudent(readyStudent());

    assertEquals(1, books.students.size());
    assertEquals(1, front.announced.size());
    assertEquals(kept.id(), front.announced.get(0).id());
  }

  @Test
  void enrollStudentRefusesWhatTheRegistryDenies() {
    front.registered = false;

    assertEquals(Errors.RUT_UNREGISTERED, assertThrows(Fault.class, () -> service.enrollStudent(readyStudent())));
    assertTrue(books.students.isEmpty());
    assertTrue(front.announced.isEmpty());
  }

  @Test
  void enrollStudentNeverTreatsAnUnknownAsAYes() {
    front.down = new IOException("socket closed");

    assertEquals(Errors.REGISTRY_UNAVAILABLE, assertThrows(Fault.class, () -> service.enrollStudent(readyStudent())));
    assertTrue(books.students.isEmpty());

    SchoolService unasked = new SchoolService(new SchoolRules(18), books, books, null, null,
        System.getLogger("school.test"));
    assertEquals(Errors.REGISTRY_UNAVAILABLE, assertThrows(Fault.class, () -> unasked.enrollStudent(readyStudent())));
  }

  @Test
  void enrollStudentKeepsAFailedNoticeInsideTheAnswer() {
    front.noticeFailure = new IOException("mail is down");

    service.enrollStudent(readyStudent());

    assertEquals(1, books.students.size());
  }

  @Test
  void theRegistryIsAskedAfterTheFormAndTheCourse() {
    Student nameless = new Student(StudentId.UNNUMBERED, new Rut("11111111-1"), new FullName(""), 20, new CourseId(1));

    assertEquals(Errors.NAME_IS_EMPTY, assertThrows(Fault.class, () -> service.enrollStudent(nameless)));
    assertEquals(Errors.COURSE_UNKNOWN,
        assertThrows(Fault.class, () -> service.enrollStudent(withCourse(readyStudent(), 9))));
    assertEquals(Errors.COURSE_UNKNOWN,
        assertThrows(Fault.class, () -> service.enrollStudent(withCourse(readyStudent(), 0))));
    assertTrue(front.asked.isEmpty());
  }

  @Test
  void theStoreIsWrittenOnlyAfterTheRegistryAnswers() {
    front.registered = false;
    assertThrows(Fault.class, () -> service.enrollStudent(readyStudent()));

    assertEquals(1, front.asked.size());
    assertTrue(books.students.isEmpty());
  }

  @Test
  void saveStudentAsksTheRegistryAndAnnouncesNothing() {
    Student kept = service.enrollStudent(readyStudent());
    int announcedByEnrolment = front.announced.size();

    front.registered = false;
    assertEquals(Errors.RUT_UNREGISTERED, assertThrows(Fault.class, () -> service.saveStudent(kept)));

    front.registered = true;
    service.saveStudent(kept);
    assertEquals(announcedByEnrolment, front.announced.size());
  }

  @Test
  void saveStudentRefusesAStudentNobodyHoldsOnlyAtTheStore() {
    Student ghost = readyStudent().withId(new StudentId(42));

    assertEquals(Errors.STUDENT_UNKNOWN, assertThrows(Fault.class, () -> service.saveStudent(ghost)));

    front.registered = false;
    assertEquals(Errors.RUT_UNREGISTERED, assertThrows(Fault.class, () -> service.saveStudent(ghost)));
  }

  @Test
  void readStudentFindsWhatEnrolmentKept() {
    Student kept = service.enrollStudent(readyStudent());

    assertEquals(kept, service.readStudent(kept.id()));
  }

  @Test
  void readStudentRefusesWhoNeverEnrolled() {
    assertEquals(Errors.STUDENT_UNKNOWN, assertThrows(Fault.class, () -> service.readStudent(new StudentId(42))));
  }

  @Test
  void listStudentsReadsTheRequestedPageOfTheRoll() {
    service.enrollStudent(readyStudent());

    assertEquals(1, service.listStudents(new Page(2, 5)).size());
    assertEquals(new Page(2, 5), books.page);
  }

  @Test
  void deleteStudentEndsAnEnrolmentOnce() {
    Student kept = service.enrollStudent(readyStudent());

    service.deleteStudent(kept.id());

    assertEquals(Errors.STUDENT_UNKNOWN, assertThrows(Fault.class, () -> service.readStudent(kept.id())));
    assertEquals(Errors.STUDENT_UNKNOWN, assertThrows(Fault.class, () -> service.deleteStudent(kept.id())));
  }

  @Test
  void createCourseOpensACourseTheCatalogueLists() {
    Course opened = service.createCourse(new Course(CourseId.UNNUMBERED, new CourseCode("FIS-201"),
        new FullName("Physics")));

    assertEquals(opened, service.readCourse(opened.id()));
    assertEquals(2, service.listCourses(new Page(0, 10)).size());
    assertEquals(10, books.page.size());
  }

  @Test
  void createCourseRefusesACourseWithoutCodeOrNameBeforeTheStoreWrites() {
    assertEquals(Errors.CODE_IS_EMPTY, assertThrows(Fault.class, () -> service.createCourse(
        new Course(CourseId.UNNUMBERED, new CourseCode(""), new FullName("Physics")))));
    assertEquals(Errors.NAME_IS_EMPTY, assertThrows(Fault.class, () -> service.createCourse(
        new Course(CourseId.UNNUMBERED, new CourseCode("FIS-201"), new FullName("")))));
    assertEquals(1, books.courses.size());
  }

  @Test
  void readCourseRefusesACourseThatNeverOpened() {
    assertEquals(Errors.COURSE_UNKNOWN, assertThrows(Fault.class, () -> service.readCourse(new CourseId(99))));
  }
}
