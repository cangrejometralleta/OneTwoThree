package com.example.app;

import com.example.faults.Fault;
import com.example.school.Page;
import com.example.transport.Request;
import com.fasterxml.jackson.databind.ObjectMapper;

import java.nio.charset.StandardCharsets;

/**
 * A Request can be Malformed without Breaking a single Business Rule.
 * Those Failures Belong to the Application.
 */
public final class Validation {

  private Validation() {
  }

  public static final Fault BODY_IS_BROKEN = Fault.refuseInvalidInput("body is not valid JSON");
  public static final Fault PATH_IS_BROKEN = Fault.refuseInvalidInput("path Holds no valid Identity");

  private static final ObjectMapper MAPPER = new ObjectMapper();

  /** ReadPathNumber Pulls an Identity out of the Path. */
  public static long readPathNumber(Request request) {
    try {
      long id = Long.parseLong(request.pathValue("id"));
      if (id < 0) {
        throw PATH_IS_BROKEN;
      }
      return id;
    } catch (NumberFormatException broken) {
      throw PATH_IS_BROKEN;
    }
  }

  /**
   * ReadJSONBody Decodes the Wire into whatever Shape the Caller Expects.
   * It Validates the Form, never the Meaning: the Core Owns the Rules.
   */
  public static <T> T readJsonBody(Request request, Class<T> shape) {
    try {
      T body = MAPPER.readValue(new String(request.body(), StandardCharsets.UTF_8), shape);
      if (body == null) {
        throw BODY_IS_BROKEN;
      }
      return body;
    } catch (Exception broken) {
      throw BODY_IS_BROKEN;
    }
  }

  /** ReadPageRequest Reads Pagination, Defaulting to the whole Set. */
  public static Page readPageRequest(Request request) {
    Page page = new Page(wholeOr(request.queryValue("page")), wholeOr(request.queryValue("size")));
    page.checkPageBounds();

    return page;
  }

  private static int wholeOr(String value) {
    try {
      return value.isBlank() ? 0 : Integer.parseInt(value);
    } catch (NumberFormatException notANumber) {
      return -1;
    }
  }
}
