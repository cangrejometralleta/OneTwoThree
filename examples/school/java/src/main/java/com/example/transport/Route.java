package com.example.transport;

/** Route Declares one Door: Method, Pattern and the Handler behind it. */
public record Route(String method, String pattern, Handler handle) {
}
