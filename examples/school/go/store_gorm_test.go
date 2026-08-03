package main

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// BuildTestingDatabase Opens SQLite in Memory, Shaped like Production.
func BuildTestingDatabase(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("database Refused to Open: %v", err)
	}

	if err := db.AutoMigrate(&StudentRow{}, &CourseRow{}); err != nil {
		t.Fatalf("schema Refused to Migrate: %v", err)
	}

	return db
}

// The Comparison that a Fake cannot Pin: a real unique Index.
func TestInsertStudentRowRefusesADuplicateRut(t *testing.T) {
	school := GormSchool{DB: BuildTestingDatabase(t)}
	student := Student{Rut: "12345678-5", Name: "Ada", Age: 20, Course: 1}

	if _, err := school.InsertStudentRow(student); err != nil {
		t.Fatalf("the first Enrolment Must Succeed, got %v", err)
	}

	if _, err := school.InsertStudentRow(student); err != ErrRutTaken {
		t.Fatalf("wanted %v, got %v", ErrRutTaken, err)
	}
}
