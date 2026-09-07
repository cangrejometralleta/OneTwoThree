package com.example.wire;

/**
 * StudentBody is what a Client Sends.
 * Weak on Purpose: the Wire Cannot be Trusted,
 * so nothing here is a Business Type yet.
 */
public record StudentBody(String rut, String name, Integer age, Long courseId) {
}
