package school

import "testing"

// Known RUTs, Checked by Hand against the Modulo eleven Rule.
func TestRutLooksValidAcceptsRealNumbers(t *testing.T) {
	valid := []RUT{"12345678-5", "11111111-1", "8459162-9", "6-k", "1-9"}

	for _, rut := range valid {
		if !rut.LooksValid() {
			t.Errorf("%s Should be valid, Got digit %q", rut, CheckDigitFor(string(rut[:len(rut)-2])))
		}
	}
}

func TestRutLooksValidRefusesBadDigits(t *testing.T) {
	invalid := []RUT{"12345678-9", "11111111-2", "12345678", "abc-1", "", "12345678-"}

	for _, rut := range invalid {
		if rut.LooksValid() {
			t.Errorf("%s Should be Refused", rut)
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

func TestCheckStudentRecordGuardsEachRule(t *testing.T) {
	good := Student{Rut: "12345678-5", Name: "Ada", Age: 20, Course: 1}

	cases := map[string]struct {
		student Student
		want    error
	}{
		"valid":       {good, nil},
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
