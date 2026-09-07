package school

import (
	"fmt"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/faults"
)

// The Business Fails in named Ways, never in Numbers.
// Each one Declares the Answer it Deserves, once and here.
var (
	ErrNameIsEmpty    = faults.RefuseInvalidInput("name is Empty")
	ErrRutIsInvalid   = faults.RefuseInvalidInput("rut Fails its Check Digit")
	ErrAgeIsTooLow    = faults.RefuseInvalidInput(fmt.Sprintf("age Must be %d or more", ReadSchoolConstants().MinimumAgeYears))
	ErrPageIsInvalid  = faults.RefuseInvalidInput("page Numbers Must not be negative")
	ErrStudentUnknown = faults.ReportMissingRecord("student not Found")
	ErrCourseUnknown  = faults.ReportMissingRecord("course not Found")
	ErrRutTaken       = faults.ReportTakenValue("rut is already Registered")
)
