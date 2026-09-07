package com.example.school;

import java.util.List;

/**
 * CourseStore Keeps Courses, and Answers whether one Exists.
 * Fulfilled by store.School.
 */
public interface CourseStore {

  Course insertCourseRow(Course course);

  Course selectCourseRow(CourseId id);

  List<Course> selectCoursePage(Page page);
}
