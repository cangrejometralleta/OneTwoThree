package store

import (
	"errors"
	"math"
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

func studentOf(rut school.RUT, name school.FullName) school.Student {
	return school.Student{Rut: rut, Name: name, Age: 20, Course: 1}
}

func TestInsertStudentRowRefusesADuplicateRut(t *testing.T) {
	books := openTestStore(t)
	if _, err := books.InsertStudentRow(studentOf("12345678-5", "Ada")); err != nil {
		t.Fatal(err)
	}

	if _, err := books.InsertStudentRow(studentOf("12345678-5", "Grace")); !errors.Is(err, school.ErrRutTaken) {
		t.Fatalf("the unique Index must Prove it, err=%v", err)
	}
}

func TestAStudentKeepsWhatItWasGivenUntrimmed(t *testing.T) {
	books := openTestStore(t)
	kept, err := books.InsertStudentRow(studentOf("12345678-5", " Ada "))
	if err != nil || kept.ID == 0 {
		t.Fatalf("kept=%v err=%v", kept, err)
	}

	if found, _, _ := books.SelectStudentRow(kept.ID); found != kept {
		t.Fatalf("found=%v kept=%v", found, kept)
	}

	kept.Name = "Ada Lovelace"
	if updated, err := books.UpdateStudentRow(kept); err != nil || !updated {
		t.Fatalf("updated=%v err=%v", updated, err)
	}
	if found, _, _ := books.SelectStudentRow(kept.ID); found.Name != "Ada Lovelace" {
		t.Fatalf("the Rewrite must Stick, found=%v", found)
	}
}

func TestTheStoreSaysNothingOfWhatItNeverKept(t *testing.T) {
	books := openTestStore(t)
	ghost := studentOf("12345678-5", "Ada")
	ghost.ID = 9

	if _, found, err := books.SelectStudentRow(9); found || err != nil {
		t.Errorf("select found=%v err=%v", found, err)
	}
	if _, found, err := books.SelectStudentRow(math.MaxInt64); found || err != nil {
		t.Errorf("select largest found=%v err=%v", found, err)
	}
	if updated, err := books.UpdateStudentRow(ghost); updated || err != nil {
		t.Errorf("update updated=%v err=%v", updated, err)
	}
	if deleted, err := books.DeleteStudentRow(9); deleted || err != nil {
		t.Errorf("delete deleted=%v err=%v", deleted, err)
	}
	if _, found, err := books.SelectCourseRow(9); found || err != nil {
		t.Errorf("course found=%v err=%v", found, err)
	}
}

func TestStudentsListByNameInBinaryOrderTiesByIdAndWholeWhenNoWindowIsAsked(t *testing.T) {
	books := openTestStore(t)
	for _, student := range []school.Student{
		studentOf("11111111-1", "carla"), studentOf("22222222-2", "Berta"),
		studentOf("12345678-5", "Ana"), studentOf("10000013-k", "Ana"),
	} {
		if _, err := books.InsertStudentRow(student); err != nil {
			t.Fatal(err)
		}
	}

	whole, _ := books.SelectStudentPage(school.WholeSet)
	order := []school.RUT{"12345678-5", "10000013-k", "22222222-2", "11111111-1"}
	for position, rut := range order {
		if whole[position].Rut != rut {
			t.Fatalf("the whole Set must List by Name in binary Order, got %v", whole)
		}
	}

	second, _ := books.SelectStudentPage(school.Page{Number: 1, Size: 3})
	if len(second) != 1 || second[0].Name != "carla" {
		t.Fatalf("the second Window of three holds the last, got %v", second)
	}
}

func TestAWindowBeyondTheRollIsEmpty(t *testing.T) {
	books := openTestStore(t)
	books.InsertStudentRow(studentOf("11111111-1", "Ana"))

	for _, page := range []school.Page{{Number: math.MaxInt64, Size: 10}, {Number: 5, Size: 10}} {
		if roll, err := books.SelectStudentPage(page); err != nil || len(roll) != 0 {
			t.Fatalf("page=%v roll=%v err=%v", page, roll, err)
		}
	}
}

func TestCoursesListByCodeAndCodesAreNotUnique(t *testing.T) {
	books := openTestStore(t)
	for _, code := range []school.CourseCode{"MAT-101", "FIS-201", "ART-100", "ART-100"} {
		if _, err := books.InsertCourseRow(school.Course{Code: code, Name: "Name"}); err != nil {
			t.Fatal(err)
		}
	}

	catalogue, err := books.SelectCoursePage(school.WholeSet)

	if err != nil || len(catalogue) != 4 || catalogue[0].Code != "ART-100" || catalogue[1].Code != "ART-100" || catalogue[3].Code != "MAT-101" {
		t.Fatalf("catalogue=%v err=%v", catalogue, err)
	}
}

func TestADeletedStudentIsGoneOnce(t *testing.T) {
	books := openTestStore(t)
	kept, _ := books.InsertStudentRow(studentOf("12345678-5", "Ada"))

	if deleted, err := books.DeleteStudentRow(kept.ID); err != nil || !deleted {
		t.Fatalf("deleted=%v err=%v", deleted, err)
	}
	if deleted, _ := books.DeleteStudentRow(kept.ID); deleted {
		t.Fatal("a second Drop must Find nobody")
	}
}
