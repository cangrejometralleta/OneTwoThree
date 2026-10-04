package school

import "testing"

func TestRutLooksValidAcceptsRealNumbers(t *testing.T) {
	for _, rut := range []RUT{"12345678-5", "11111111-1", "22222222-2", "10000013-k", "10000004-0"} {
		if !rut.LooksValid() {
			t.Errorf("%s Should be Accepted", rut)
		}
	}
}

func TestRutLooksValidRefusesBadDigits(t *testing.T) {
	for _, rut := range []RUT{"12345678-9", "12345678", "abc-1", "-5", "12345678-55", ""} {
		if rut.LooksValid() {
			t.Errorf("%s Should be Refused", rut)
		}
	}
}

func TestCheckStudentRecordGuardsEachRule(t *testing.T) {
	cases := map[string]struct {
		student Student
		want    error
	}{
		"valid":       {Student{Rut: "12345678-5", Name: "Ada", Age: 20, Course: 1}, nil},
		"minimum age": {Student{Rut: "12345678-5", Name: "Ada", Age: 18}, nil},
		"no name":     {Student{Rut: "12345678-5", Age: 20}, ErrNameIsEmpty},
		"bad rut":     {Student{Rut: "12345678-9", Name: "Ada", Age: 20}, ErrRutIsInvalid},
		"too young":   {Student{Rut: "12345678-5", Name: "Ada", Age: 17}, ErrAgeIsTooLow},
	}

	for name, test := range cases {
		if got := test.student.CheckStudentRecord(); got != test.want {
			t.Errorf("%s: wanted %v, got %v", name, test.want, got)
		}
	}
}

func TestCheckCourseRecordGuardsEachRule(t *testing.T) {
	cases := map[string]struct {
		course Course
		want   error
	}{
		"valid":   {Course{Code: "MAT-101", Name: "Algebra"}, nil},
		"no code": {Course{Name: "Algebra"}, ErrCodeIsEmpty},
		"no name": {Course{Code: "MAT-101"}, ErrNameIsEmpty},
		"neither": {Course{}, ErrCodeIsEmpty},
	}

	for name, test := range cases {
		if got := test.course.CheckCourseRecord(); got != test.want {
			t.Errorf("%s: wanted %v, got %v", name, test.want, got)
		}
	}
}

func TestReadSchoolConstantsHoldsTheEnrolmentAge(t *testing.T) {
	if got := ReadSchoolConstants().MinimumAgeYears; got != 18 {
		t.Fatalf("the Enrolment Age is eighteen, got %d", got)
	}
}

func TestCheckPageBoundsRefusesWhatTheStoreCannotServe(t *testing.T) {
	if (Page{Number: 0, Size: 0}).CheckPageBounds() != nil || (Page{Number: 3, Size: 10}).CheckPageBounds() != nil {
		t.Fatal("a Window and the whole Set are both Servable")
	}
	if (Page{Number: -1, Size: 10}).CheckPageBounds() != ErrPageIsInvalid || (Page{Number: 0, Size: -5}).CheckPageBounds() != ErrPageIsInvalid {
		t.Fatal("a Page Counting backwards must be Refused")
	}
}
