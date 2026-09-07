package com.example.transport;

/** Response is what a Handler Returns. */
public record Response(int status, Object body) {
}
