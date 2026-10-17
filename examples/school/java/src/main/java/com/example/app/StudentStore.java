package com.example.app;

import com.example.school.Page;
import com.example.school.Student;
import com.example.school.StudentId;

import java.util.List;
import java.util.Optional;

/**
 * StudentStore Keeps Students wherever Students Live.
 * Absence is an empty Optional or false; only the Store Proves a Rut Taken.
 */
public interface StudentStore {
  Student insertStudentRow(Student student);

  Optional<Student> selectStudentRow(StudentId id);

  boolean updateStudentRow(Student student);

  boolean deleteStudentRow(StudentId id);

  List<Student> selectStudentPage(Page page);
}
