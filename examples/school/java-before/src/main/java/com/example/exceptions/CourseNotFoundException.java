package com.example.exceptions;

import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.http.HttpStatus;

@ResponseStatus(code = HttpStatus.NOT_FOUND)
public class CourseNotFoundException extends Exception {
  private static final long serialVersionUID = -9205427303291243642L;

  /**
   * This exception must be thrown when a single course is searched and it is not
   * found in it's repository
   * 
   * @param courseId The id of the course that was not found
   */
  public CourseNotFoundException(Long courseId) {
    super("The course with id: '" + courseId + "', was not found");
  }
}