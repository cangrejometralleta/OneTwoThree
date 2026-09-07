package com.example.services;

import com.example.exceptions.CourseNotFoundException;
import com.example.exceptions.InvalidPaginationParameterException;
import com.example.entities.Course;
import com.example.repositories.CourseRepository;

import java.util.Optional;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.PageRequest;
import org.springframework.stereotype.Service;

@Service
public class CourseService {

  @Autowired
  private CourseRepository courseRepository;

  /**
   * 
   * @return all the courses found in the database
   */
  public Iterable<Course> findAll() {
    return courseRepository.findAll();
  }

  /**
   * 
   * @param pageNumber
   * @param pageSize
   * @return all the courses obtained in the current page
   */
  public Iterable<Course> findByPage(int pageNumber, int pageSize) throws InvalidPaginationParameterException {
    if (pageNumber < 0 || pageSize < 0) {
      throw new InvalidPaginationParameterException();
    }

    Pageable pageable = PageRequest.of(pageNumber, pageSize);
    return courseRepository.findAll(pageable);
  }

  /**
   * 
   * @param courseId the courseId of the course to find
   * @return the course found in the repository
   * @throws CourseNotFoundException
   */
  public Course findById(Long courseId) throws CourseNotFoundException {
    Optional<Course> course = courseRepository.findById(courseId);
    if (!course.isPresent()) {
      throw new CourseNotFoundException(courseId);
    }
    return course.get();
  }

  /**
   * 
   * @param course the course to save in the database
   */
  public Course create(Course course) {
    return courseRepository.save(course);
  }

  public Course update(Course course) throws CourseNotFoundException {
    Course updateCourse = this.findById(course.getId());
    updateCourse.setCode(course.getCode());
    updateCourse.setName(course.getName());

    return courseRepository.save(updateCourse);
  }

  /**
   * 
   * @param courseId
   * @throws CourseNotFoundException
   */
  public void delete(Long courseId) throws CourseNotFoundException {
    Course course = this.findById(courseId);
    courseRepository.delete(course);
  }
}