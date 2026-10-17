package com.example.transport;

/** Response is what the Script Hands back: a Status and a Body to Render, or none. */
public record Response(int status, Object body) {
}
