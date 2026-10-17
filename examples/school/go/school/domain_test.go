package school

import (
	"errors"
	"math"
	"testing"
)

var rules = SchoolRules{MinimumAgeYears: 18}

func TestRutLooksValidAcceptsRealNumbers(t *testing.T) {
	for _, rut := range []RUT{"12345678-5", "11111111-1", "22222222-2", "10000013-k", "10000013-K", "10000004-0",
		"123456789012345678901234567890-1"} {
		if !rut.LooksValid() {
			t.Errorf("%s Should be Accepted", rut)
		}
	}
}

func TestRutLooksValidRefusesBadDigits(t *testing.T) {
	for _, rut := range []RUT{"12345678-9", "12345678", "abc-1", "-5", "12345678-55", "", " 12345678-5"} {
		if rut.LooksValid() {
			t.Errorf("%s Should be Refused", rut)
		}
	}
}

func TestARutEntersInLowerCaseAndKeepsItsLeadingZeros(t *testing.T) {
	if ReadRUT("10000013-K") != "10000013-k" {
		t.Fatalf("got %s", ReadRUT("10000013-K"))
	}
	if ReadRUT("012345678-5") == ReadRUT("12345678-5") {
		t.Fatal("leading zeros must not Merge")
	}
}

func TestCheckStudentRecordGuardsEachRule(t *testing.T) {
	cases := map[string]struct {
		student Student
		want    error
	}{
		"valid":              {Student{Rut: "12345678-5", Name: "Ada", Age: 20, Course: 1}, nil},
		"minimum age":        {Student{Rut: "12345678-5", Name: "Ada", Age: 18}, nil},
		"no name first":      {Student{Rut: "bad", Age: 1}, ErrNameIsEmpty},
		"only whitespace":    {Student{Rut: "12345678-5", Name: " \t", Age: 20}, ErrNameIsEmpty},
		"bad rut before age": {Student{Rut: "12345678-9", Name: "Ada", Age: 1}, ErrRutIsInvalid},
		"too young":          {Student{Rut: "12345678-5", Name: "Ada", Age: 17}, ErrAgeIsTooLow(rules)},
	}

	for name, test := range cases {
		if got := test.student.CheckStudentRecord(rules); !errors.Is(got, test.want) {
			t.Errorf("%s: wanted %v, got %v", name, test.want, got)
		}
	}
}

func TestTheAgeReasonSpellsTheFloorFromTheConstant(t *testing.T) {
	if got := ErrAgeIsTooLow(rules).Error(); got != "age Must be 18 or more" {
		t.Fatalf("got %q", got)
	}
	if got := ErrAgeIsTooLow(SchoolRules{MinimumAgeYears: 21}).Error(); got != "age Must be 21 or more" {
		t.Fatalf("got %q", got)
	}
}

func TestCheckCourseRecordGuardsEachRule(t *testing.T) {
	cases := map[string]struct {
		course Course
		want   error
	}{
		"valid":   {Course{Code: "MAT-101", Name: "Algebra"}, nil},
		"no code": {Course{Name: "Algebra"}, ErrCodeIsEmpty},
		"no name": {Course{Code: "MAT-101", Name: " "}, ErrNameIsEmpty},
		"neither": {Course{}, ErrCodeIsEmpty},
	}

	for name, test := range cases {
		if got := test.course.CheckCourseRecord(); got != test.want {
			t.Errorf("%s: wanted %v, got %v", name, test.want, got)
		}
	}
}

func TestAPageRefusesWhatTheStoreCannotServe(t *testing.T) {
	if _, err := BuildPage(0, 1); err != nil {
		t.Fatal(err)
	}
	if _, err := BuildPage(-1, 10); err != ErrPageIsInvalid {
		t.Fatalf("err=%v", err)
	}
	if _, err := BuildPage(0, 0); err != ErrPageIsInvalid {
		t.Fatalf("err=%v", err)
	}
}

func TestAPageBeyondSixtyFourBitsWindowsAsFarAsItCan(t *testing.T) {
	if got := (Page{Number: 2, Size: 10}).Offset(); got != 20 {
		t.Fatalf("got %d", got)
	}
	if got := (Page{Number: math.MaxInt64, Size: 10}).Offset(); got != math.MaxInt64 {
		t.Fatalf("got %d", got)
	}
}
