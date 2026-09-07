package com.example.controllers;

import com.example.exceptions.CourseNotFoundException;
import com.example.exceptions.InvalidPaginationParameterException;
import com.example.entities.Course;
import com.example.services.CourseService;

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
@RequestMapping("/courses")
public class CourseController {

  @Autowired
  private CourseService courseService;

  @GetMapping("/")
  public Iterable<Course> findByPage(int pageNumber, int pageSize) throws InvalidPaginationParameterException {
    return courseService.findByPage(pageNumber, pageSize);
  }

  @GetMapping("/all")
  public Iterable<Course> findAll() {
    return courseService.findAll();
  }

  @GetMapping("/{courseId}")
  public Course get(@PathVariable Long courseId) throws CourseNotFoundException {
    return courseService.findById(courseId);
  }

  @ResponseStatus(code = HttpStatus.CREATED)
  @PostMapping("/")
  public Course create(@Valid @RequestBody Course course) {
    return courseService.create(course);
  }

  @PutMapping("/{courseId}")
  public Course update(@PathVariable Long courseId, @Valid @RequestBody Course course) throws CourseNotFoundException {
    course.setId(courseId);
    return courseService.update(course);
  }

  @DeleteMapping("/{courseId}")
  public void delete(@PathVariable Long courseId) throws CourseNotFoundException {
    courseService.delete(courseId);
  }
}