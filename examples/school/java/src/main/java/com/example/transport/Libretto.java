package com.example.transport;

import java.util.List;

/**
 * Libretto is every Route a Service Declares, and the one Answer for a Door it does not Hold.
 * The Unknown Answers before any Guard, because there is no Route to Guard.
 */
public record Libretto(List<Route> routes, Handler unknown) {

  public Libretto {
    routes = List.copyOf(routes);
  }
}
