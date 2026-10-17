package school

import (
	"fmt"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/faults"
)

// The Business Fails in named Ways, never in Numbers.
// Each one Declares the Answer it Deserves, once and here.
var (
	ErrNameIsEmpty         = faults.RefuseInvalidInput("name is Empty")
	ErrCodeIsEmpty         = faults.RefuseInvalidInput("code is Empty")
	ErrRutIsInvalid        = faults.RefuseInvalidInput("rut Fails its Check Digit")
	ErrPageIsInvalid       = faults.RefuseInvalidInput("page or size is not valid")
	ErrStudentUnknown      = faults.ReportMissingRecord("student not Found")
	ErrCourseUnknown       = faults.ReportMissingRecord("course not Found")
	ErrRutTaken            = faults.ReportTakenValue("rut is already Registered")
	ErrRutUnregistered     = faults.RefuseUnprocessableValue("rut is not in the Registry")
	ErrRegistryUnavailable = faults.ReportUnavailableProvider("registry cannot Answer now")
)

// ErrAgeIsTooLow Spells the Floor from the Global Constant.
func ErrAgeIsTooLow(rules SchoolRules) error {
	return faults.RefuseInvalidInput(fmt.Sprintf("age Must be %d or more", rules.MinimumAgeYears))
}
