package com.example.exceptions;

import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.http.HttpStatus;

@ResponseStatus(code = HttpStatus.NOT_FOUND)
public class StudentNotFoundException extends Exception {

  /**
   *
   */
  private static final long serialVersionUID = -5630940129955610309L;

  /**
   * This exception must be thrown when a single student is searched and it is not
   * found in it's repository
   * 
   * @param courseId The id of the student that was not found
   */
  public StudentNotFoundException(Long studentId) {
    super("The student with id: '" + studentId + "' was not found");
  }
}