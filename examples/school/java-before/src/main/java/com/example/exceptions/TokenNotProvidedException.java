package com.example.exceptions;

import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.http.HttpStatus;

@ResponseStatus(code = HttpStatus.UNAUTHORIZED)
public class TokenNotProvidedException extends Exception {

  private static final long serialVersionUID = 5094909385121600903L;

  /**
   * This exception is used when a token was not provided to a request
   */
  public TokenNotProvidedException() {
    super("A token was not provided with this request");
  }
}