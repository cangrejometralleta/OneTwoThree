package com.example.transport;

/** Server Serves a Libretto on a Port; each Framework Fills it in its own Package. */
public interface Server {
  Running serveRoutes(Libretto libretto, int port) throws Exception;

  /** Running is a Server that Listens until it is Closed. */
  interface Running extends AutoCloseable {
    int port();

    @Override
    void close();
  }
}
