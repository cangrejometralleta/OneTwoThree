package com.example.store;

import com.example.school.Course;
import com.example.school.Errors;
import com.example.school.Page;
import com.example.school.CourseStore;
import com.example.school.StudentStore;
import com.example.school.Student;
import com.example.school.CourseCode;
import com.example.school.CourseId;
import com.example.school.FullName;
import com.example.school.Rut;
import com.example.school.StudentId;
import com.example.store.Rows.CourseRow;
import com.example.store.Rows.StudentRow;

import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Sort;

import java.util.List;

/**
 * School Fulfils both Stores with one Connection.
 * It is the only Type in this Program that Knows JPA Exists.
 */
public record School(StudentRows studentRows, CourseRows courseRows) implements StudentStore, CourseStore {

  /** InsertStudentRow Writes a new Student and Returns it Numbered. */
  @Override
  public Student insertStudentRow(Student student) {
    try {
      return decodeStudentRow(studentRows.saveAndFlush(encodeStudentRow(student, null)));
    } catch (DataIntegrityViolationException collision) {
      throw Errors.RUT_TAKEN;
    }
  }

  /** SelectStudentRow Finds one Student or Says why not. */
  @Override
  public Student selectStudentRow(StudentId id) {
    return decodeStudentRow(studentRows.findById(id.value()).orElseThrow(() -> Errors.STUDENT_UNKNOWN));
  }

  /** UpdateStudentRow Overwrites an existing Student. */
  @Override
  public Student updateStudentRow(Student student) {
    StudentRow row = studentRows.findById(student.id().value()).orElseThrow(() -> Errors.STUDENT_UNKNOWN);

    try {
      return decodeStudentRow(studentRows.saveAndFlush(encodeStudentRow(student, row.id)));
    } catch (DataIntegrityViolationException collision) {
      throw Errors.RUT_TAKEN;
    }
  }

  /** DeleteStudentRow Removes a Student, or Reports an Absence. */
  @Override
  public void deleteStudentRow(StudentId id) {
    if (!studentRows.existsById(id.value())) {
      throw Errors.STUDENT_UNKNOWN;
    }

    studentRows.deleteById(id.value());
  }

  /** SelectStudentPage Reads one Page, Ordered by Name. */
  @Override
  public List<Student> selectStudentPage(Page page) {
    Sort byName = Sort.by("name");
    List<StudentRow> rows = page.size() <= 0
        ? studentRows.findAll(byName)
        : studentRows.findAll(PageRequest.of(page.number(), page.size(), byName)).getContent();

    return rows.stream().map(School::decodeStudentRow).toList();
  }

  /** InsertCourseRow Writes a new Course and Returns it Numbered. */
  @Override
  public Course insertCourseRow(Course course) {
    CourseRow row = new CourseRow();
    row.code = course.code().value();
    row.name = course.name().value();

    return decodeCourseRow(courseRows.saveAndFlush(row));
  }

  /** SelectCourseRow Finds one Course or Says why not. */
  @Override
  public Course selectCourseRow(CourseId id) {
    return decodeCourseRow(courseRows.findById(id.value()).orElseThrow(() -> Errors.COURSE_UNKNOWN));
  }

  /** SelectCoursePage Reads one Page of Courses, Ordered by Code. */
  @Override
  public List<Course> selectCoursePage(Page page) {
    Sort byCode = Sort.by("code");
    List<CourseRow> rows = page.size() <= 0
        ? courseRows.findAll(byCode)
        : courseRows.findAll(PageRequest.of(page.number(), page.size(), byCode)).getContent();

    return rows.stream().map(School::decodeCourseRow).toList();
  }

  /** EncodeStudentRow Turns a Business Student into Storage. */
  private static StudentRow encodeStudentRow(Student student, Long id) {
    StudentRow row = new StudentRow();
    row.id = id;
    row.rut = student.rut().value();
    row.name = student.name().value();
    row.age = student.age();
    row.courseId = student.course().value();

    return row;
  }

  /** DecodeStudentRow Turns Storage back into Business. */
  private static Student decodeStudentRow(StudentRow row) {
    return new Student(new StudentId(row.id), new Rut(row.rut), new FullName(row.name), row.age,
        new CourseId(row.courseId));
  }

  /** DecodeCourseRow Turns a Course Row back into Business. */
  private static Course decodeCourseRow(CourseRow row) {
    return new Course(new CourseId(row.id), new CourseCode(row.code), new FullName(row.name));
  }
}
