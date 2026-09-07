package school

import "testing"

// A Deployment Value never Reaches a Global Constant.
func TestReadSchoolConstantsHoldsTheEnrolmentAge(t *testing.T) {
	if ReadSchoolConstants().MinimumAgeYears != 18 {
		t.Fatal("Deployment Must Preserve Global Constants")
	}
}
