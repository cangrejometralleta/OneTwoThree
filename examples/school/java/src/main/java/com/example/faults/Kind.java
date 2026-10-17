package com.example.faults;

/**
 * A Kind Names what Went wrong, in the Business' own Words.
 * It Names no Protocol: app Translates a Kind into a Status.
 */
public enum Kind {
  /** UNEXPECTED is every Failure the Program did not Plan for. */
  UNEXPECTED,
  INVALID_INPUT,
  UNPROVEN_CALLER,
  MISSING_RECORD,
  TAKEN_VALUE,
  UNPROCESSABLE_VALUE,
  UNAVAILABLE_PROVIDER
}
