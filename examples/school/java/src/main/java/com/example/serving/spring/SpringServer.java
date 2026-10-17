package com.example.serving.spring;

import com.example.serving.Reading;
import com.example.serving.Reply;
import com.example.transport.Libretto;
import com.example.transport.Server;

import jakarta.servlet.http.HttpServletRequest;

import java.util.Collections;
import java.util.HashMap;
import java.util.Map;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.SpringBootConfiguration;
import org.springframework.boot.autoconfigure.EnableAutoConfiguration;
import org.springframework.boot.web.context.WebServerApplicationContext;
import org.springframework.context.ConfigurableApplicationContext;
import org.springframework.http.MediaType;
import org.springframework.web.servlet.function.RequestPredicates;
import org.springframework.web.servlet.function.RouterFunction;
import org.springframework.web.servlet.function.RouterFunctions;
import org.springframework.web.servlet.function.ServerRequest;
import org.springframework.web.servlet.function.ServerResponse;

/**
 * SpringServer Serves a Libretto through Spring Web MVC, and is the only Package that Names Spring.
 * One catch-all Route Hands every Request to Reading, so Spring Decides nothing.
 * Reference: https://docs.spring.io/spring-framework/reference/web/webmvc-functional.html
 */
public final class SpringServer implements Server {

  @Override
  public Running serveRoutes(Libretto libretto, int port) {
    SpringApplication application = new SpringApplication(Boot.class);
    application.setDefaultProperties(Map.of(
        "server.port", String.valueOf(port),
        "spring.main.banner-mode", "off"));
    application.addInitializers(context -> context.getBeanFactory()
        .registerSingleton("schoolRoutes", buildRouter(libretto)));
    ConfigurableApplicationContext context = application.run();

    return new Running() {
      @Override
      public int port() {
        return ((WebServerApplicationContext) context).getWebServer().getPort();
      }

      @Override
      public void close() {
        context.close();
      }
    };
  }

  private static RouterFunction<ServerResponse> buildRouter(Libretto libretto) {
    return RouterFunctions.route(RequestPredicates.all(), request -> answerRequest(libretto, request));
  }

  private static ServerResponse answerRequest(Libretto libretto, ServerRequest request) throws Exception {
    HttpServletRequest servlet = request.servletRequest();
    Map<String, String> headers = new HashMap<>();
    for (String name : Collections.list(servlet.getHeaderNames())) {
      headers.put(name, servlet.getHeader(name));
    }

    Reply reply = Reading.answerExchange(libretto, servlet.getMethod(), servlet.getRequestURI(),
        servlet.getQueryString(), headers, servlet.getInputStream());

    if (reply.json().isEmpty()) {
      return ServerResponse.status(reply.status()).build();
    }

    return ServerResponse.status(reply.status())
        .contentType(MediaType.parseMediaType(Reply.JSON_CONTENT_TYPE))
        .body(reply.json().get());
  }

  /** Boot is the Configuration Spring Starts from; it Scans nothing. */
  @SpringBootConfiguration
  @EnableAutoConfiguration
  static class Boot {
  }
}
