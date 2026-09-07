package com.example.exceptions;

import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.http.HttpStatus;

@ResponseStatus(code = HttpStatus.BAD_REQUEST)
public class InvalidPaginationParameterException extends Exception {
  private static final long serialVersionUID = 2066657991909390125L;

  /**
   * This exception is used when wrong pagination values are send
   */
  public InvalidPaginationParameterException() {
    super("Pagination parameters should be positive integers or zero");
  }
}