package com.example.app;

import com.example.school.Rut;
import com.example.school.Student;

import java.util.ArrayList;
import java.util.List;

/** FrontDesk Answers both Questions, the Way one Adapter may Fill two Ports. */
public final class FrontDesk implements RutRegistry, EnrollmentNotifier {

  public boolean registered = true;
  public Exception down;
  public Exception noticeFailure;
  public final List<Rut> asked = new ArrayList<>();
  public final List<Student> announced = new ArrayList<>();

  @Override
  public boolean confirmRut(Rut rut) throws Exception {
    asked.add(rut);
    if (down != null) {
      throw down;
    }

    return registered;
  }

  @Override
  public void announceEnrollment(Student student) throws Exception {
    announced.add(student);
    if (noticeFailure != null) {
      throw noticeFailure;
    }
  }
}
