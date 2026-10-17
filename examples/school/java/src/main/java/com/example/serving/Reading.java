package com.example.serving;

import static java.net.HttpURLConnection.HTTP_NO_CONTENT;

import com.example.codec.Json;
import com.example.transport.Libretto;
import com.example.transport.Request;
import com.example.transport.Response;
import com.example.transport.Route;

import java.io.IOException;
import java.io.InputStream;
import java.net.URLDecoder;
import java.nio.charset.StandardCharsets;
import java.util.HashMap;
import java.util.Locale;
import java.util.Map;
import java.util.Optional;

/**
 * Reading Joins the three Beats every Adapter Shares: read, answer, write.
 * An Adapter Hands it Bytes, a Method, a raw Path and Headers; it Hands back a Reply.
 * So every Framework Answers the same Way, because none of them Decides.
 */
public final class Reading {

  private Reading() {
  }

  /** answerExchange Finds the Route, Reads the Request and Renders the Answer. */
  public static Reply answerExchange(Libretto libretto, String method, String rawPath, String rawQuery,
      Map<String, String> headers, InputStream body) throws IOException {
    Map<String, String> lowered = new HashMap<>();
    headers.forEach((name, value) -> lowered.putIfAbsent(name.toLowerCase(Locale.ROOT), value));
    Request request = new Request(Map.of(), readQueryValues(rawQuery), lowered,
        body.readNBytes(Request.MAX_BODY_BYTES + 1), "");

    for (Route route : libretto.routes()) {
      Optional<Map<String, String>> names = matchRoutePattern(route.pattern(), rawPath);
      if (route.method().equals(method) && names.isPresent()) {
        return writeReply(route.handle().handle(request.withPath(names.get())));
      }
    }

    return writeReply(libretto.unknown().handle(request));
  }

  /** matchRoutePattern Reads the Names a Pattern Declares: /students/{id} Holds id. /students/ is not /students. */
  static Optional<Map<String, String>> matchRoutePattern(String pattern, String rawPath) {
    String[] wanted = pattern.split("/", -1);
    String[] given = rawPath.split("/", -1);
    if (wanted.length != given.length) {
      return Optional.empty();
    }

    Map<String, String> names = new HashMap<>();
    for (int position = 0; position < wanted.length; position++) {
      String part = wanted[position];
      if (part.startsWith("{") && part.endsWith("}")) {
        Optional<String> value = decodePathSegment(given[position]);
        if (given[position].isEmpty() || value.isEmpty()) {
          return Optional.empty();
        }
        names.put(part.substring(1, part.length() - 1), value.get());
      } else if (!part.equals(given[position])) {
        return Optional.empty();
      }
    }

    return Optional.of(names);
  }

  /** readQueryValues Takes the first Value of a repeated Key; a Pair that cannot be Decoded is Dropped. */
  static Map<String, String> readQueryValues(String rawQuery) {
    Map<String, String> values = new HashMap<>();
    if (rawQuery == null || rawQuery.isEmpty()) {
      return values;
    }

    for (String pair : rawQuery.split("&")) {
      int equals = pair.indexOf('=');
      String name = equals < 0 ? pair : pair.substring(0, equals);
      String value = equals < 0 ? "" : pair.substring(equals + 1);
      try {
        values.putIfAbsent(URLDecoder.decode(name, StandardCharsets.UTF_8),
            URLDecoder.decode(value, StandardCharsets.UTF_8));
      } catch (IllegalArgumentException undecodable) {
        continue;
      }
    }

    return values;
  }

  /** writeReply Renders the Script's Response: JSON, and no Body for a 204. */
  static Reply writeReply(Response response) {
    if (response.status() == HTTP_NO_CONTENT || response.body() == null) {
      return new Reply(response.status(), Optional.empty());
    }

    return new Reply(response.status(),
        Optional.of(Json.writeJsonText(response.body()).getBytes(StandardCharsets.UTF_8)));
  }

  private static Optional<String> decodePathSegment(String segment) {
    try {
      return Optional.of(URLDecoder.decode(segment.replace("+", "%2B"), StandardCharsets.UTF_8));
    } catch (IllegalArgumentException undecodable) {
      return Optional.empty();
    }
  }
}
