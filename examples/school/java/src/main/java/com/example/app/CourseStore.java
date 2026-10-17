package com.example.app;

import com.example.school.Course;
import com.example.school.CourseId;
import com.example.school.Page;

import java.util.List;
import java.util.Optional;

/** CourseStore Keeps Courses, and Answers whether one Exists. */
public interface CourseStore {
  Course insertCourseRow(Course course);

  Optional<Course> selectCourseRow(CourseId id);

  List<Course> selectCoursePage(Page page);
}
