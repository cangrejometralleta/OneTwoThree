package com.example.store;

/** StudentRow is the Storage Shape of a Student: it Knows its Table and no one else does. */
record StudentRow(long id, String rut, String name, long age, long courseId) {
}
