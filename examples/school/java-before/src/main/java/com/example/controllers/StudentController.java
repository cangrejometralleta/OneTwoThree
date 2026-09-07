package com.example.controllers;

import com.example.exceptions.StudentNotFoundException;
import com.example.exceptions.CourseNotFoundException;
import com.example.exceptions.InvalidPaginationParameterException;
import com.example.entities.Student;
import com.example.services.StudentService;

import javax.validation.Valid;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.bind.annotation.ResponseStatus;

@RestController
@RequestMapping("/students")
public class StudentController {

  @Autowired
  private StudentService studentService;

  @GetMapping("/")
  public Iterable<Student> findByPage(int pageNumber, int pageSize) throws InvalidPaginationParameterException {
    return studentService.findByPage(pageNumber, pageSize);
  }

  @GetMapping("/all")
  public Iterable<Student> findAll() {
    return studentService.findAll();
  }

  @GetMapping("/{studentId}")
  public Student get(@PathVariable Long studentId) throws StudentNotFoundException {
    return studentService.findById(studentId);
  }

  @ResponseStatus(code = HttpStatus.CREATED)
  @PostMapping("/")
  public Student create(@RequestBody @Valid Student student) throws CourseNotFoundException {
    return studentService.create(student);
  }

  @PutMapping("/{studentId}")
  public Student update(@PathVariable Long studentId, @RequestBody @Valid Student student)
      throws CourseNotFoundException, StudentNotFoundException {
    student.setId(studentId);
    return studentService.update(student);
  }

  @DeleteMapping("/{studentId}")
  public void delete(@PathVariable Long studentId) throws StudentNotFoundException {
    studentService.delete(studentId);
  }
}