package com.example.store;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.Table;

/** The Storage Shapes. The Tags Document the Table. */
public final class Rows {

  private Rows() {
  }

  @Entity
  @Table(name = "student")
  public static class StudentRow {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(name = "student_id")
    public Long id;

    @Column(name = "rut", length = 16, unique = true, nullable = false)
    public String rut;

    @Column(name = "name", length = 120, nullable = false)
    public String name;

    @Column(name = "age", nullable = false)
    public int age;

    @Column(name = "course_id", nullable = false)
    public long courseId;
  }

  @Entity
  @Table(name = "course")
  public static class CourseRow {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(name = "course_id")
    public Long id;

    @Column(name = "code", length = 16, unique = true, nullable = false)
    public String code;

    @Column(name = "name", length = 120, nullable = false)
    public String name;
  }
}
