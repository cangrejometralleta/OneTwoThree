package com.example.store;

import com.example.store.Rows.CourseRow;

import org.springframework.data.jpa.repository.JpaRepository;

/** CourseRows is the Spring Data Door, and only this Package Opens it. */
public interface CourseRows extends JpaRepository<CourseRow, Long> {
}
