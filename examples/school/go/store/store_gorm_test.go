package store

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/school"
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
	store := School{DB: BuildTestingDatabase(t)}
	student := school.Student{Rut: "12345678-5", Name: "Ada", Age: 20, Course: 1}

	if _, err := store.InsertStudentRow(student); err != nil {
		t.Fatalf("the first Enrolment Must Succeed, got %v", err)
	}

	if _, err := store.InsertStudentRow(student); err != school.ErrRutTaken {
		t.Fatalf("wanted %v, got %v", school.ErrRutTaken, err)
	}
}
