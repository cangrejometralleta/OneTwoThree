package com.example.wire;

/** CourseView is what a Client Gets back about a Course. */
public record CourseView(long id, String code, String name) {
}
