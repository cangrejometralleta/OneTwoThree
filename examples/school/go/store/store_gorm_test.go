package store

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/app"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/school"
)

// One Store Fills both Ports, and the Compiler Checks it.
var (
	_ app.StudentStore = School{}
	_ app.CourseStore  = School{}
)

func openTestStore(t *testing.T) School {
	t.Helper()
	opened, err := OpenSchoolStore(filepath.Join(t.TempDir(), "school.db"))
	if err != nil {
		t.Fatal(err)
	}

	return opened
}

func TestInsertStudentRowRefusesADuplicateRut(t *testing.T) {
	books := openTestStore(t)
	ada := school.Student{Rut: "12345678-5", Name: "Ada", Age: 20, Course: 1}

	if _, err := books.InsertStudentRow(ada); err != nil {
		t.Fatal(err)
	}

	if _, err := books.InsertStudentRow(ada); !errors.Is(err, school.ErrRutTaken) {
		t.Fatalf("the unique Index must Prove it, err=%v", err)
	}
}

func TestAStudentKeepsWhatItWasGiven(t *testing.T) {
	books := openTestStore(t)
	kept, err := books.InsertStudentRow(school.Student{Rut: "12345678-5", Name: "Ada", Age: 20, Course: 1})
	if err != nil || kept.ID == 0 {
		t.Fatalf("kept=%v err=%v", kept, err)
	}

	found, err := books.SelectStudentRow(kept.ID)
	if err != nil || found != kept {
		t.Fatalf("found=%v kept=%v err=%v", found, kept, err)
	}

	kept.Name = "Ada Lovelace"
	if _, err := books.UpdateStudentRow(kept); err != nil {
		t.Fatal(err)
	}
	if found, _ = books.SelectStudentRow(kept.ID); found.Name != "Ada Lovelace" {
		t.Fatalf("the Rewrite must Stick, found=%v", found)
	}
}

func TestTheStoreSaysUnknownForWhatItNeverKept(t *testing.T) {
	books := openTestStore(t)

	if _, err := books.SelectStudentRow(9); !errors.Is(err, school.ErrStudentUnknown) {
		t.Errorf("select err=%v", err)
	}
	if _, err := books.UpdateStudentRow(school.Student{ID: 9, Rut: "12345678-5", Name: "Ada", Age: 20, Course: 1}); !errors.Is(err, school.ErrStudentUnknown) {
		t.Errorf("update err=%v", err)
	}
	if err := books.DeleteStudentRow(9); !errors.Is(err, school.ErrStudentUnknown) {
		t.Errorf("delete err=%v", err)
	}
	if _, err := books.SelectCourseRow(9); !errors.Is(err, school.ErrCourseUnknown) {
		t.Errorf("course err=%v", err)
	}
}

func TestStudentsListByNameInWindowsAndWholeWhenNoWindowIsAsked(t *testing.T) {
	books := openTestStore(t)
	for index, name := range []string{"Carla", "Ana", "Berta"} {
		rut := []school.RUT{"11111111-1", "22222222-2", "12345678-5"}[index]
		if _, err := books.InsertStudentRow(school.Student{Rut: rut, Name: school.FullName(name), Age: 20, Course: 1}); err != nil {
			t.Fatal(err)
		}
	}

	whole, _ := books.SelectStudentPage(school.Page{})
	second, _ := books.SelectStudentPage(school.Page{Number: 1, Size: 2})

	if len(whole) != 3 || whole[0].Name != "Ana" || whole[2].Name != "Carla" {
		t.Fatalf("the whole Set must List by Name, got %v", whole)
	}
	if len(second) != 1 || second[0].Name != "Carla" {
		t.Fatalf("the second Window of two holds the last, got %v", second)
	}
}

func TestCoursesListByCode(t *testing.T) {
	books := openTestStore(t)
	for _, code := range []school.CourseCode{"MAT-101", "FIS-201", "ART-100"} {
		if _, err := books.InsertCourseRow(school.Course{Code: code, Name: "Name"}); err != nil {
			t.Fatal(err)
		}
	}

	catalogue, err := books.SelectCoursePage(school.Page{})

	if err != nil || len(catalogue) != 3 || catalogue[0].Code != "ART-100" || catalogue[2].Code != "MAT-101" {
		t.Fatalf("catalogue=%v err=%v", catalogue, err)
	}
}

func TestADeletedStudentIsGoneOnce(t *testing.T) {
	books := openTestStore(t)
	kept, _ := books.InsertStudentRow(school.Student{Rut: "12345678-5", Name: "Ada", Age: 20, Course: 1})

	if err := books.DeleteStudentRow(kept.ID); err != nil {
		t.Fatal(err)
	}
	if err := books.DeleteStudentRow(kept.ID); !errors.Is(err, school.ErrStudentUnknown) {
		t.Fatalf("a second Drop must be Unknown, err=%v", err)
	}
}
