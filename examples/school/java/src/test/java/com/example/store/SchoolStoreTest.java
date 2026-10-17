package com.example.store;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
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
import com.example.school.Student;
import com.example.school.StudentId;

import java.nio.file.Path;
import java.sql.SQLException;
import java.util.List;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;

class SchoolStoreTest {

  @TempDir
  Path folder;

  private SchoolStore books;

  @BeforeEach
  void openTestStore() throws SQLException {
    books = SchoolStore.openSchoolStore(folder.resolve("school"));
  }

  private static Student studentOf(String rut, String name) {
    return new Student(StudentId.UNNUMBERED, new Rut(rut), new FullName(name), 20, new CourseId(1));
  }

  @Test
  void theUniqueIndexProvesARutIsTaken() {
    books.insertStudentRow(studentOf("12345678-5", "Ada"));

    assertEquals(Errors.RUT_TAKEN,
        assertThrows(Fault.class, () -> books.insertStudentRow(studentOf("12345678-5", "Grace"))));
  }

  @Test
  void aStudentKeepsWhatItWasGivenUntrimmed() {
    Student kept = books.insertStudentRow(studentOf("12345678-5", " Ada "));

    assertTrue(kept.id().value() > 0);
    assertEquals(kept, books.selectStudentRow(kept.id()).orElseThrow());

    Student renamed = new Student(kept.id(), kept.rut(), new FullName("Ada Lovelace"), 21, kept.course());
    assertTrue(books.updateStudentRow(renamed));
    assertEquals(renamed, books.selectStudentRow(kept.id()).orElseThrow());
  }

  @Test
  void theStoreSaysNothingOfWhatItNeverKept() {
    assertTrue(books.selectStudentRow(new StudentId(9)).isEmpty());
    assertFalse(books.updateStudentRow(studentOf("12345678-5", "Ada").withId(new StudentId(9))));
    assertFalse(books.deleteStudentRow(new StudentId(9)));
    assertTrue(books.selectCourseRow(new CourseId(9)).isEmpty());
    assertTrue(books.selectStudentRow(new StudentId(Long.MAX_VALUE)).isEmpty());
  }

  @Test
  void studentsListByNameInBinaryOrderTiesByIdAndWholeWhenNoWindowIsAsked() {
    books.insertStudentRow(studentOf("11111111-1", "carla"));
    books.insertStudentRow(studentOf("22222222-2", "Berta"));
    books.insertStudentRow(studentOf("12345678-5", "Ana"));
    books.insertStudentRow(studentOf("10000013-k", "Ana"));

    List<String> whole = books.selectStudentPage(Page.WHOLE_SET).stream().map(s -> s.rut().text()).toList();
    assertEquals(List.of("12345678-5", "10000013-k", "22222222-2", "11111111-1"), whole);

    List<Student> second = books.selectStudentPage(new Page(1, 3));
    assertEquals(1, second.size());
    assertEquals("carla", second.get(0).name().text());
  }

  @Test
  void aWindowBeyondTheRollIsEmpty() {
    books.insertStudentRow(studentOf("11111111-1", "Ana"));

    assertEquals(List.of(), books.selectStudentPage(new Page(Long.MAX_VALUE, 10)));
    assertEquals(List.of(), books.selectStudentPage(new Page(5, 10)));
  }

  @Test
  void coursesListByCodeAndCodesAreNotUnique() {
    for (String code : List.of("MAT-101", "FIS-201", "ART-100", "ART-100")) {
      books.insertCourseRow(new Course(CourseId.UNNUMBERED, new CourseCode(code), new FullName("Name")));
    }

    List<String> codes = books.selectCoursePage(Page.WHOLE_SET).stream().map(c -> c.code().text()).toList();
    assertEquals(List.of("ART-100", "ART-100", "FIS-201", "MAT-101"), codes);
  }

  @Test
  void aDeletedStudentIsGoneOnce() {
    Student kept = books.insertStudentRow(studentOf("12345678-5", "Ada"));

    assertTrue(books.deleteStudentRow(kept.id()));
    assertFalse(books.deleteStudentRow(kept.id()));
  }

  @Test
  void theStoreSurvivesBeingOpenedAgain() throws SQLException {
    Student kept = books.insertStudentRow(studentOf("12345678-5", "Ada"));

    SchoolStore reopened = SchoolStore.openSchoolStore(folder.resolve("school"));

    assertEquals(kept, reopened.selectStudentRow(kept.id()).orElseThrow());
  }
}
