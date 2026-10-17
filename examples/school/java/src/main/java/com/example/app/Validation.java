package com.example.app;

import com.example.codec.Json;
import com.example.faults.Fault;
import com.example.school.Errors;
import com.example.school.Page;
import com.example.transport.Request;

import java.math.BigInteger;
import java.util.Locale;
import java.util.Map;
import java.util.Optional;
import java.util.Set;
import java.util.regex.Pattern;

/** Validation Reads the Form of a Request, before its Meaning: the Path, the Body, the Page. */
public final class Validation {

  public static final Fault BODY_IS_BROKEN = Fault.refuseInvalidInput("body is not valid JSON");
  public static final Fault PATH_IS_BROKEN = Fault.refuseInvalidInput("path Holds no valid Identity");
  public static final Fault ROUTE_UNKNOWN = Fault.reportMissingRecord("route not Found");

  private static final Pattern DIGITS = Pattern.compile("^[0-9]+$");
  private static final BigInteger LARGEST_WHOLE = BigInteger.valueOf(Long.MAX_VALUE);
  private static final String BEARER_SCHEME = "bearer ";

  private Validation() {
  }

  /**
   * readPathNumber Reads the positive whole Number in the Path's {id}; 012 is 12.
   * A Number beyond 64 bits Names nobody, so it Reads as the largest, which no one Holds.
   */
  public static long readPathNumber(Request request) {
    String text = request.path().getOrDefault("id", "");
    if (!DIGITS.matcher(text).matches()) {
      throw PATH_IS_BROKEN;
    }

    long number = readClampedNumber(text);
    if (number == 0) {
      throw PATH_IS_BROKEN;
    }

    return number;
  }

  /**
   * readJsonObject Reads one JSON Object whose Fields are exactly those Named, none null.
   * Anything Less, or More, is Body is Broken.
   */
  public static Map<String, Object> readJsonObject(Request request, Set<String> fields) {
    byte[] body = request.body();
    if (body.length > Request.MAX_BODY_BYTES) {
      throw BODY_IS_BROKEN;
    }

    Optional<Object> value = Json.readJsonText(body);
    if (value.isEmpty() || !(value.get() instanceof Map<?, ?> sent) || !sent.keySet().equals(fields)
        || sent.containsValue(Json.NULL)) {
      throw BODY_IS_BROKEN;
    }

    @SuppressWarnings("unchecked")
    Map<String, Object> object = (Map<String, Object>) sent;

    return object;
  }

  /** readTextField Reads a Field that Must be a String. */
  public static String readTextField(Map<String, Object> object, String name) {
    if (object.get(name) instanceof String text) {
      return text;
    }

    throw BODY_IS_BROKEN;
  }

  /** readWholeField Reads a Field that Must be a whole Number of 64 bits. */
  public static long readWholeField(Map<String, Object> object, String name) {
    if (object.get(name) instanceof BigInteger whole && whole.bitLength() < Long.SIZE) {
      return whole.longValueExact();
    }

    throw BODY_IS_BROKEN;
  }

  /**
   * readPageRequest Reads the Window a list Route Asks for.
   * Neither Named Asks for the whole Set; one alone Takes the other's Default.
   */
  public static Page readPageRequest(Request request) {
    String page = request.query().get("page");
    String size = request.query().get("size");
    if (page == null && size == null) {
      return Page.WHOLE_SET;
    }

    return new Page(page == null ? 0 : readPageNumber(page),
        size == null ? Page.DEFAULT_PAGE_SIZE : readPageNumber(size));
  }

  /** readBearerToken Reads the Token after a Bearer Scheme of any Case, or Nothing. */
  public static String readBearerToken(Request request) {
    String header = request.headers().getOrDefault("authorization", "");
    if (header.length() < BEARER_SCHEME.length()
        || !header.substring(0, BEARER_SCHEME.length()).toLowerCase(Locale.ROOT).equals(BEARER_SCHEME)) {
      return "";
    }

    return header.substring(BEARER_SCHEME.length());
  }

  private static long readPageNumber(String text) {
    if (!DIGITS.matcher(text).matches()) {
      throw Errors.PAGE_IS_INVALID;
    }

    return readClampedNumber(text);
  }

  private static long readClampedNumber(String digits) {
    return new BigInteger(digits).min(LARGEST_WHOLE).longValueExact();
  }
}
