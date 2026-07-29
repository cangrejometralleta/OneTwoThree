package main

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

func TestCheckStudentRecordGuardsEachRule(t *testing.T) {
	good := Student{Rut: "12345678-5", Name: "Ada", Age: 20, Course: 1}

	cases := map[string]struct {
		student Student
		want    error
	}{
		"valid":     {good, nil},
		"no name":   {Student{Rut: "12345678-5", Age: 20}, ErrNameIsEmpty},
		"bad rut":   {Student{Rut: "12345678-9", Name: "Ada", Age: 20}, ErrRutIsInvalid},
		"too young": {Student{Rut: "12345678-5", Name: "Ada", Age: 17}, ErrAgeIsTooLow},
	}

	for name, test := range cases {
		if got := test.student.CheckStudentRecord(); got != test.want {
			t.Errorf("%s: wanted %v, got %v", name, test.want, got)
		}
	}
}
