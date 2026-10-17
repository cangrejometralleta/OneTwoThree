package com.example.school;

/**
 * Page Windows a List, Zero based. A Page of WHOLE_SET Asks for every Record.
 * The Bounds Hold before a Page Exists: Number at least zero, Size at least one.
 */
public record Page(long number, long size) {

  /** DEFAULT_PAGE_SIZE Windows a Page when the Caller Names a page but no size. */
  public static final long DEFAULT_PAGE_SIZE = 10;

  /** WHOLE_SET Stands for a List no Window Cuts. */
  public static final Page WHOLE_SET = new Page(0, Long.MAX_VALUE);

  public Page {
    if (number < 0 || size < 1) {
      throw Errors.PAGE_IS_INVALID;
    }
  }

  /** offset Counts the Records before the Window, as far as 64 bits can Count. */
  public long offset() {
    long high = Math.multiplyHigh(number, size);
    long low = number * size;

    return high != 0 || low < 0 ? Long.MAX_VALUE : low;
  }
}
