package com.example.entities;

import com.example.validation.RUTConstraint;

import javax.persistence.Column;
import javax.persistence.Entity;
import javax.persistence.GenerationType;
import javax.persistence.GeneratedValue;
import javax.persistence.Id;
import javax.persistence.JoinColumn;
import javax.persistence.ManyToOne;
import javax.persistence.Table;
import javax.validation.constraints.Min;
import javax.validation.constraints.NotBlank;
import javax.validation.constraints.NotNull;

@Entity
@Table(name = "student")
public class Student {
  @Column(name = "student_id")
  @GeneratedValue(strategy = GenerationType.AUTO)
  @Id
  private Long id;

  @Column(name = "rut")
  @NotBlank(message = "RUT is mandatory")
  @NotNull
  @RUTConstraint
  private String rut;
  @Column(name = "name")
  @NotBlank(message = "Name is mandatory")
  @NotNull

  private String name;
  @Column(name = "last_name")
  @NotBlank(message = "LastName is mandatory")
  @NotNull

  private String lastName;
  @Column(name = "age")
  @Min(value = 18, message = "Age should not be less than 18")
  @NotNull
  private int age;

  @ManyToOne
  @NotNull
  @JoinColumn(name = "course_id")
  private Course course;

  public Student() {
  }

  public Student(Long id, String rut, String name, String lastName, int age, Course course) {
    this.id = id;
    this.rut = rut;
    this.name = name;
    this.lastName = lastName;
    this.age = age;
    this.course = course;
  }

  public Long getId() {
    return this.id;
  }

  public void setId(Long id) {
    this.id = id;
  }

  public String getRut() {
    return this.rut;
  }

  public void setRut(String rut) {
    this.rut = rut;
  }

  public String getName() {
    return this.name;
  }

  public void setName(String name) {
    this.name = name;
  }

  public String getLastName() {
    return this.lastName;
  }

  public void setLastName(String lastName) {
    this.lastName = lastName;
  }

  public int getAge() {
    return this.age;
  }

  public void setAge(int age) {
    this.age = age;
  }

  public Course getCourse() {
    return this.course;
  }

  public void setCourse(Course course) {
    this.course = course;
  }

  @Override
  public String toString() {
    return "{" + " id='" + getId() + "'" + ", rut='" + getRut() + "'" + ", name='" + getName() + "'" + ", lastName='"
        + getLastName() + "'" + ", age='" + getAge() + "'" + ", course='" + getCourse() + "'" + "}";
  }

}