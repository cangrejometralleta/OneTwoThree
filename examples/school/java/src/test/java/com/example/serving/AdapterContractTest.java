package com.example.serving;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

import com.example.serving.spring.SpringServer;
import com.example.serving.stdlib.StdlibServer;
import com.example.transport.Libretto;
import com.example.transport.Request;
import com.example.transport.Response;
import com.example.transport.Route;
import com.example.transport.Server;

import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.List;
import java.util.Map;
import java.util.concurrent.atomic.AtomicReference;
import java.util.stream.Stream;

import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.MethodSource;

/**
 * Every Framework is Held to one Contract: the same Path, Query, Headers and Body Arrive,
 * and the same Status and JSON Leave.
 */
class AdapterContractTest {

  private static final HttpClient CLIENT = HttpClient.newHttpClient();

  /** ANY_FREE_PORT Asks the System for a Port nobody is Using. */
  private static final int ANY_FREE_PORT = 0;

  static Stream<Server> everyAdapter() {
    return Stream.of(new StdlibServer(), new SpringServer());
  }

  private static HttpResponse<String> send(Server.Running running, String method, String target, String body,
      Map<String, String> headers) throws Exception {
    HttpRequest.Builder request = HttpRequest.newBuilder(URI.create("http://localhost:" + running.port() + target))
        .method(method, body == null ? HttpRequest.BodyPublishers.noBody() : HttpRequest.BodyPublishers.ofString(body));
    headers.forEach(request::header);

    return CLIENT.send(request.build(), HttpResponse.BodyHandlers.ofString());
  }

  @ParameterizedTest
  @MethodSource("everyAdapter")
  void everyAdapterAnswersTheSameWay(Server adapter) throws Exception {
    AtomicReference<Request> seen = new AtomicReference<>();
    Libretto libretto = new Libretto(List.of(
        new Route("POST", "/echo/{id}", request -> {
          seen.set(request);

          return new Response(201, Map.of("id", request.path().get("id")));
        }),
        new Route("GET", "/echo", request -> new Response(200, List.of())),
        new Route("DELETE", "/echo/{id}", request -> new Response(204, null)),
        new Route("GET", "/failing", request -> new Response(404, Map.of("error", "gone")))),
        request -> new Response(404, Map.of("error", "route not Found")));

    try (Server.Running running = adapter.serveRoutes(libretto, ANY_FREE_PORT)) {
      HttpResponse<String> echoed = send(running, "POST", "/echo/007?page=2&size=5&page=9", "{\"a\":1}",
          Map.of("Authorization", "bEaReR abc", "Content-Type", "text/plain"));
      assertEquals(201, echoed.statusCode());
      assertEquals("application/json", echoed.headers().firstValue("Content-Type").orElse("").split(";")[0]);
      assertEquals("{\"id\":\"007\"}", echoed.body());
      assertEquals("007", seen.get().path().get("id"));
      assertEquals(Map.of("page", "2", "size", "5"), seen.get().query());
      assertEquals("bEaReR abc", seen.get().headers().get("authorization"));
      assertEquals("{\"a\":1}", new String(seen.get().body(), java.nio.charset.StandardCharsets.UTF_8));

      assertEquals("[]", send(running, "GET", "/echo", null, Map.of()).body());

      HttpResponse<String> dropped = send(running, "DELETE", "/echo/7", null, Map.of());
      assertEquals(204, dropped.statusCode());
      assertTrue(dropped.body().isEmpty());

      assertEquals(404, send(running, "GET", "/failing", null, Map.of()).statusCode());

      for (String[] stranger : new String[][] {{"GET", "/nowhere"}, {"GET", "/echo/"}, {"PUT", "/echo/7"},
          {"GET", "/echo/7/more"}, {"PATCH", "/echo"}}) {
        HttpResponse<String> unknown = send(running, stranger[0], stranger[1], null, Map.of());
        assertEquals(404, unknown.statusCode(), stranger[0] + " " + stranger[1]);
        assertEquals("{\"error\":\"route not Found\"}", unknown.body(), stranger[0] + " " + stranger[1]);
      }
    }
  }

  @ParameterizedTest
  @MethodSource("everyAdapter")
  void aBodyBeyondTheBoundArrivesOneBytePastIt(Server adapter) throws Exception {
    AtomicReference<Request> seen = new AtomicReference<>();
    Libretto libretto = new Libretto(List.of(new Route("POST", "/echo", request -> {
      seen.set(request);

      return new Response(204, null);
    })), request -> new Response(404, null));

    try (Server.Running running = adapter.serveRoutes(libretto, ANY_FREE_PORT)) {
      send(running, "POST", "/echo", " ".repeat(Request.MAX_BODY_BYTES * 2), Map.of());

      assertEquals(Request.MAX_BODY_BYTES + 1, seen.get().body().length);
    }
  }
}
