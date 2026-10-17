package com.example.app;

import com.example.school.Course;
import com.example.school.CourseCode;
import com.example.school.CourseId;
import com.example.school.Errors;
import com.example.school.FullName;
import com.example.school.Page;
import com.example.school.Student;
import com.example.school.StudentId;

import java.util.ArrayList;
import java.util.Comparator;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.TreeMap;

/**
 * Shelf Stands in for the Store: Course 1 is Open, and it Keeps what it is Given.
 * It Obeys the Store's Rules, so a Rut is Unique here too.
 */
public final class Shelf implements StudentStore, CourseStore {

  public final Map<Long, Student> students = new TreeMap<>();
  public final Map<Long, Course> courses = new TreeMap<>(Map.of(1L,
      new Course(new CourseId(1), new CourseCode("MAT-101"), new FullName("Algebra"))));
  public Page page;
  private long nextStudent = 1;

  @Override
  public Student insertStudentRow(Student student) {
    if (students.values().stream().anyMatch(kept -> kept.rut().equals(student.rut()))) {
      throw Errors.RUT_TAKEN;
    }
    Student kept = student.withId(new StudentId(nextStudent++));
    students.put(kept.id().value(), kept);

    return kept;
  }

  @Override
  public Optional<Student> selectStudentRow(StudentId id) {
    return Optional.ofNullable(students.get(id.value()));
  }

  @Override
  public boolean updateStudentRow(Student student) {
    return students.replace(student.id().value(), student) != null;
  }

  @Override
  public boolean deleteStudentRow(StudentId id) {
    return students.remove(id.value()) != null;
  }

  @Override
  public List<Student> selectStudentPage(Page asked) {
    page = asked;
    List<Student> roll = new ArrayList<>(students.values());
    roll.sort(Comparator.comparing((Student s) -> s.name().text()).thenComparing(s -> s.id().value()));

    return roll;
  }

  @Override
  public Course insertCourseRow(Course course) {
    Course kept = course.withId(new CourseId(courses.size() + 1));
    courses.put(kept.id().value(), kept);

    return kept;
  }

  @Override
  public Optional<Course> selectCourseRow(CourseId id) {
    return Optional.ofNullable(courses.get(id.value()));
  }

  @Override
  public List<Course> selectCoursePage(Page asked) {
    page = asked;

    return new ArrayList<>(courses.values());
  }
}
