package com.example.transport;

/**
 * Route Binds one Method and one Pattern to one Handler.
 * Patterns Use {name}, and each Adapter Translates from there.
 */
public record Route(String method, String pattern, Handler handle) {
}
