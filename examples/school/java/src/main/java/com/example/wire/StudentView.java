package com.example.wire;

/** StudentView is what a Client Gets back. */
public record StudentView(long id, String rut, String name, int age, long courseId) {
}
