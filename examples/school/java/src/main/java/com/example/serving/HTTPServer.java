package com.example.serving;

import com.example.transport.Request;
import com.example.transport.Response;
import com.example.transport.Route;

import org.springframework.http.MediaType;
import org.springframework.web.servlet.function.RouterFunction;
import org.springframework.web.servlet.function.RouterFunctions;
import org.springframework.web.servlet.function.ServerRequest;
import org.springframework.web.servlet.function.ServerResponse;

import java.nio.charset.StandardCharsets;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

/**
 * HTTPServer Mounts a Route Table on whatever Web Framework it Wraps.
 * The Name Says the Responsibility; only the Imports Say Spring,
 * and this is the only Type in the Program that Names them.
 * Reference: https://docs.spring.io/spring-framework/reference/web/webmvc-functional.html
 */
public final class HTTPServer {

  private final List<Route> routes;

  /**
   * The Routes Arrive through the Constructor, Cast in Main.
   * The Adapter Discovers nothing and Scans nothing.
   */
  public HTTPServer(List<Route> routes) {
    this.routes = List.copyOf(routes);
  }

  /** ServeRoutes Mounts every Route on one Router. */
  public RouterFunction<ServerResponse> serveRoutes() {
    RouterFunctions.Builder builder = RouterFunctions.route();

    for (Route route : routes) {
      builder.route(
          serverRequest -> serverRequest.method().name().equals(route.method())
              && serverRequest.path().equals(pathOf(serverRequest, route)),
          serverRequest -> writeReplyAsJson(route.handle().handle(readRequestValues(serverRequest, route))));
    }

    return builder.build();
  }

  /** ReadRequestValues Copies a Spring Request into our own Shape. */
  private static Request readRequestValues(ServerRequest serverRequest, Route route) {
    Map<String, String> query = new LinkedHashMap<>();
    serverRequest.params().forEach((key, values) -> query.put(key, values.get(0)));

    return new Request(pathParamsOf(serverRequest, route), query, readBearerToken(serverRequest),
        readBodyBytes(serverRequest), "");
  }

  /** ReadBodyBytes Hands the Body over as Bytes, whatever Spring Threw Reading it. */
  private static byte[] readBodyBytes(ServerRequest serverRequest) {
    try {
      return serverRequest.body(String.class).getBytes(StandardCharsets.UTF_8);
    } catch (Exception unreadable) {
      return new byte[0];
    }
  }

  /** ReadBearerToken Pulls the Token out of the Authorization Header. */
  private static String readBearerToken(ServerRequest serverRequest) {
    String header = serverRequest.headers().firstHeader("Authorization");

    return header == null ? "" : header.replaceFirst("^Bearer ", "");
  }

  /** WriteReplyAsJSON is the one Place that Touches a ServerResponse. */
  private static ServerResponse writeReplyAsJson(Response reply) {
    ServerResponse.BodyBuilder builder = ServerResponse.status(reply.status()).contentType(MediaType.APPLICATION_JSON);

    return reply.body() == null ? builder.build() : builder.body(reply.body());
  }

  private static Map<String, String> pathParamsOf(ServerRequest serverRequest, Route route) {
    String[] wanted = route.pattern().split("/");
    String[] given = serverRequest.path().split("/");

    Map<String, String> path = new LinkedHashMap<>();
    for (int index = 0; index < wanted.length && index < given.length; index++) {
      if (wanted[index].startsWith("{") && wanted[index].endsWith("}")) {
        path.put(wanted[index].substring(1, wanted[index].length() - 1), given[index]);
      }
    }

    return path;
  }

  private static String pathOf(ServerRequest serverRequest, Route route) {
    String[] wanted = route.pattern().split("/");
    String[] given = serverRequest.path().split("/");
    if (wanted.length != given.length) {
      return "";
    }

    StringBuilder rebuilt = new StringBuilder();
    for (int index = 1; index < wanted.length; index++) {
      boolean parameter = wanted[index].startsWith("{") && wanted[index].endsWith("}");
      rebuilt.append('/').append(parameter ? given[index] : wanted[index]);
    }

    return rebuilt.toString();
  }
}
