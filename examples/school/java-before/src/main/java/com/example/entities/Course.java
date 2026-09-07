package com.example.entities;

import javax.persistence.Column;
import javax.persistence.Entity;
import javax.persistence.GenerationType;
import javax.persistence.GeneratedValue;
import javax.persistence.Id;
import javax.persistence.Table;
import javax.validation.constraints.NotBlank;
import javax.validation.constraints.NotNull;
import javax.validation.constraints.Size;

@Entity
@Table(name = "course")
public class Course {
  @Column(name = "course_id")
  @GeneratedValue(strategy = GenerationType.AUTO)
  @Id
  private Long id;

  @Column(name = "name")
  @NotBlank(message = "Name is mandatory")
  @NotNull
  private String name;

  @Column(name = "code")
  @NotBlank(message = "Code is mandatory")
  @NotNull
  @Size(max = 4, message = "Code must have a maximum length of 4")
  private String code;

  public Course() {
  }

  public Course(Long id, String name, String code) {
    this.id = id;
    this.name = name;
    this.code = code;
  }

  public Long getId() {
    return this.id;
  }

  public void setId(Long id) {
    this.id = id;
  }

  public String getName() {
    return this.name;
  }

  public void setName(String name) {
    this.name = name;
  }

  public String getCode() {
    return this.code;
  }

  public void setCode(String code) {
    this.code = code;
  }

  @Override
  public String toString() {
    return "{" + " id='" + getId() + "'" + ", name='" + getName() + "'" + ", code='" + getCode() + "'" + "}";
  }

}