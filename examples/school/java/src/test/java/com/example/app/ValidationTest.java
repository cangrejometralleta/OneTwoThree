package com.example.app;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

import com.example.faults.Fault;
import com.example.school.Errors;
import com.example.school.Page;
import com.example.transport.Request;

import java.nio.charset.StandardCharsets;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

import org.junit.jupiter.api.Test;

class ValidationTest {

  private static final Set<String> STUDENT_FIELDS = Set.of("rut", "name", "age", "courseId");
  private static final String WHOLE = "{\"rut\":\"11111111-1\",\"name\":\"Ana\",\"age\":20,\"courseId\":1}";

  private static Request pathOf(String id) {
    return new Request(Map.of("id", id), Map.of(), Map.of(), new byte[0], "");
  }

  private static Request bodyOf(byte[] body) {
    return new Request(Map.of(), Map.of(), Map.of(), body, "");
  }

  private static Request queryOf(Map<String, String> query) {
    return new Request(Map.of(), query, Map.of(), new byte[0], "");
  }

  private static Map<String, Object> readStudentFields(byte[] body) {
    Map<String, Object> sent = Validation.readJsonObject(bodyOf(body), STUDENT_FIELDS);
    Validation.readTextField(sent, "rut");
    Validation.readTextField(sent, "name");
    Validation.readWholeField(sent, "age");
    Validation.readWholeField(sent, "courseId");

    return sent;
  }

  @Test
  void readPathNumberWantsAPositiveWholeNumber() {
    assertEquals(7, Validation.readPathNumber(pathOf("7")));
    assertEquals(12, Validation.readPathNumber(pathOf("012")));

    for (String bad : List.of("", "abc", "0", "000", "-3", "1.5", "+1", " 1", "1e1", "0x1")) {
      assertEquals(Validation.PATH_IS_BROKEN, assertThrows(Fault.class, () -> Validation.readPathNumber(pathOf(bad))),
          bad);
    }
  }

  @Test
  void aPathNumberBeyondSixtyFourBitsNamesNobodyRatherThanBreaking() {
    assertEquals(Long.MAX_VALUE, Validation.readPathNumber(pathOf("99999999999999999999999")));
  }

  @Test
  void readJsonObjectRefusesAnythingButOneWholeObject() {
    readStudentFields(WHOLE.getBytes(StandardCharsets.UTF_8));

    Map<String, String> broken = new HashMap<>();
    broken.put("not json", "{\"rut\":");
    broken.put("empty", "");
    broken.put("missing field", "{\"rut\":\"11111111-1\",\"name\":\"Ana\",\"age\":20}");
    broken.put("null field", "{\"rut\":\"11111111-1\",\"name\":null,\"age\":20,\"courseId\":1}");
    broken.put("unknown field", "{\"rut\":\"11111111-1\",\"name\":\"Ana\",\"age\":20,\"courseId\":1,\"id\":9}");
    broken.put("wrong type", "{\"rut\":\"11111111-1\",\"name\":\"Ana\",\"age\":\"20\",\"courseId\":1}");
    broken.put("text as number", "{\"rut\":11111111,\"name\":\"Ana\",\"age\":20,\"courseId\":1}");
    broken.put("fractional", "{\"rut\":\"11111111-1\",\"name\":\"Ana\",\"age\":20.5,\"courseId\":1}");
    broken.put("whole written as float", "{\"rut\":\"11111111-1\",\"name\":\"Ana\",\"age\":20.0,\"courseId\":1}");
    broken.put("beyond 64 bits", "{\"rut\":\"11111111-1\",\"name\":\"Ana\",\"age\":20,\"courseId\":9223372036854775808}");
    broken.put("second value", WHOLE + " {}");
    broken.put("an array", "[]");
    broken.put("null document", "null");
    broken.put("repeated key", "{\"rut\":\"11111111-1\",\"rut\":\"x\",\"name\":\"Ana\",\"age\":20,\"courseId\":1}");

    broken.forEach((name, body) -> assertEquals(Validation.BODY_IS_BROKEN,
        assertThrows(Fault.class, () -> readStudentFields(body.getBytes(StandardCharsets.UTF_8))), name));
  }

  @Test
  void aBodyMustBeUtf8AndAtMostOneMebibyte() {
    byte[] latin = WHOLE.replace("Ana", "Aná").getBytes(StandardCharsets.ISO_8859_1);
    assertEquals(Validation.BODY_IS_BROKEN, assertThrows(Fault.class, () -> readStudentFields(latin)));

    byte[] oversized = new byte[Request.MAX_BODY_BYTES + 1];
    java.util.Arrays.fill(oversized, (byte) ' ');
    byte[] whole = WHOLE.getBytes(StandardCharsets.UTF_8);
    System.arraycopy(whole, 0, oversized, 0, whole.length);
    assertEquals(Validation.BODY_IS_BROKEN, assertThrows(Fault.class, () -> readStudentFields(oversized)));
  }

  @Test
  void readPageRequestWindowsOnlyWhatWasAsked() {
    assertEquals(Page.WHOLE_SET, Validation.readPageRequest(queryOf(Map.of())));
    assertEquals(new Page(2, 5), Validation.readPageRequest(queryOf(Map.of("page", "2", "size", "5"))));
    assertEquals(new Page(3, Page.DEFAULT_PAGE_SIZE), Validation.readPageRequest(queryOf(Map.of("page", "3"))));
    assertEquals(new Page(0, 25), Validation.readPageRequest(queryOf(Map.of("size", "25"))));
    assertEquals(Page.WHOLE_SET, Validation.readPageRequest(queryOf(Map.of("sort", "name"))));
  }

  @Test
  void aPageBeyondSixtyFourBitsWindowsAsFarAsItCan() {
    assertEquals(new Page(Long.MAX_VALUE, 10),
        Validation.readPageRequest(queryOf(Map.of("page", "99999999999999999999", "size", "10"))));
  }

  @Test
  void readPageRequestRefusesWhatIsNotDigitsInRange() {
    for (Map<String, String> query : List.of(Map.of("page", "abc"), Map.of("page", "-1"), Map.of("page", ""),
        Map.of("page", "1.5"), Map.of("page", "+1"), Map.of("page", "1e1"), Map.of("page", " 1"),
        Map.of("page", "0x1"), Map.of("size", "many"), Map.of("size", "0"), Map.of("size", "-4"),
        Map.of("size", ""))) {
      assertEquals(Errors.PAGE_IS_INVALID,
          assertThrows(Fault.class, () -> Validation.readPageRequest(queryOf(query))), query.toString());
    }
  }

  @Test
  void aBearerTokenIsReadOnlyAfterTheScheme() {
    assertEquals("abc", Validation.readBearerToken(
        new Request(Map.of(), Map.of(), Map.of("authorization", "bEaReR abc"), new byte[0], "")));
    assertEquals("", Validation.readBearerToken(
        new Request(Map.of(), Map.of(), Map.of("authorization", "Basic abc"), new byte[0], "")));
    assertEquals("", Validation.readBearerToken(Request.EMPTY));
  }
}
