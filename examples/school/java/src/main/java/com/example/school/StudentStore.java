package com.example.school;

import java.util.List;

/**
 * A Provider is an Interface the Core Declares
 * and something outside Fulfils.
 * The Core Depends on the Shape, never on the Library behind it.
 *
 * <p>The Word Collides. Spring Calls a registered Bean a Provider,
 * and JSR-330 Calls a Factory one. Here it Means the Door an outside
 * System Enters through, which the Literature Calls a Port.
 *
 * <p>StudentStore Keeps Students wherever Students Live.
 * Fulfilled by store.School.
 */
public interface StudentStore {

  Student insertStudentRow(Student student);

  Student selectStudentRow(StudentId id);

  Student updateStudentRow(Student student);

  void deleteStudentRow(StudentId id);

  List<Student> selectStudentPage(Page page);
}
