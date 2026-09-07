package com.example.school;

/** Page Carries Pagination without Naming a Database. */
public record Page(int number, int size) {

  /** CheckPageBounds Refuses a Page the Store cannot Serve. */
  public void checkPageBounds() {
    if (number < 0 || size < 0) {
      throw Errors.PAGE_IS_INVALID;
    }
  }
}
