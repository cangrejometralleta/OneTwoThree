package com.example.store;

import com.example.store.Rows.StudentRow;

import org.springframework.data.jpa.repository.JpaRepository;

/** StudentRows is the Spring Data Door, and only this Package Opens it. */
public interface StudentRows extends JpaRepository<StudentRow, Long> {
}
