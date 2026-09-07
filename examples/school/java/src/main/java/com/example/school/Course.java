package com.example.school;


/** Course is the Business Truth about one Class. */
public record Course(CourseId id, CourseCode code, FullName name) {
}
