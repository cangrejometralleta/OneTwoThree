package com.example.wire;

/**
 * StudentBody is what a Client Sends. Weak on Purpose: the Wire cannot be
 * Trusted, so nothing here is a Business Type yet.
 */
public record StudentBody(String rut, String name, long age, long courseId) {
}
