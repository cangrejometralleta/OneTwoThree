package com.example.app;

import com.example.school.Student;

/** EnrollmentNotifier Announces a Student the Store has Kept. */
public interface EnrollmentNotifier {
  void announceEnrollment(Student student) throws Exception;
}
