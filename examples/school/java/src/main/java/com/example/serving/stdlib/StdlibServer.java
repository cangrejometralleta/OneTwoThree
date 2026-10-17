package com.example.serving.stdlib;

import com.example.serving.Reading;
import com.example.serving.Reply;
import com.example.transport.Libretto;
import com.example.transport.Server;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.Executors;

/**
 * StdlibServer Serves a Libretto with the JDK's own HttpServer, and nothing else.
 * Reference: https://docs.oracle.com/en/java/javase/21/docs/api/jdk.httpserver/com/sun/net/httpserver/HttpServer.html
 */
public final class StdlibServer implements Server {

  /** NO_BODY Tells the JDK a Reply Carries no Bytes at all. */
  private static final int NO_BODY = -1;

  /** MAX_SWALLOW_BYTES Bounds what is Read and Dropped past the Body, as Tomcat's own maxSwallowSize does. */
  private static final long MAX_SWALLOW_BYTES = 2L << 20;
  private static final int SWALLOW_BUFFER_BYTES = 8 << 10;

  /** IMMEDIATE_STOP Closes without Waiting for open Exchanges. */
  private static final int IMMEDIATE_STOP = 0;

  @Override
  public Running serveRoutes(Libretto libretto, int port) throws IOException {
    HttpServer server = HttpServer.create(new InetSocketAddress(port), 0);
    server.createContext("/", exchange -> answerExchange(libretto, exchange));
    server.setExecutor(Executors.newVirtualThreadPerTaskExecutor());
    server.start();

    return new Running() {
      @Override
      public int port() {
        return server.getAddress().getPort();
      }

      @Override
      public void close() {
        server.stop(IMMEDIATE_STOP);
      }
    };
  }

  /** swallowUnreadBody Drops what the Bound Left, so the Client can Read the Answer before the Socket Closes. */
  private static void swallowUnreadBody(HttpExchange exchange) throws IOException {
    InputStream body = exchange.getRequestBody();
    byte[] buffer = new byte[SWALLOW_BUFFER_BYTES];
    long dropped = 0;
    int read;
    while (dropped < MAX_SWALLOW_BYTES
        && (read = body.read(buffer, 0, (int) Math.min(buffer.length, MAX_SWALLOW_BYTES - dropped))) >= 0) {
      dropped += read;
    }
  }

  private static void answerExchange(Libretto libretto, HttpExchange exchange) throws IOException {
    try (exchange) {
      Map<String, String> headers = new HashMap<>();
      exchange.getRequestHeaders().forEach((name, values) -> headers.put(name, values.get(0)));

      Reply reply = Reading.answerExchange(libretto, exchange.getRequestMethod(),
          exchange.getRequestURI().getRawPath(), exchange.getRequestURI().getRawQuery(), headers,
          exchange.getRequestBody());
      swallowUnreadBody(exchange);

      if (reply.json().isEmpty()) {
        exchange.sendResponseHeaders(reply.status(), NO_BODY);
        return;
      }

      byte[] json = reply.json().get();
      exchange.getResponseHeaders().set("Content-Type", Reply.JSON_CONTENT_TYPE);
      exchange.sendResponseHeaders(reply.status(), json.length);
      try (OutputStream out = exchange.getResponseBody()) {
        out.write(json);
      }
    }
  }
}
