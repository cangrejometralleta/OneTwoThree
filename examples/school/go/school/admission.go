package school

import (
	"path/filepath"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/settings"
)

// OldestPlausibleAge Bounds the Constant File, not the Business.
// A School Refuses at eighteen; this only Refuses a Typo.
const OldestPlausibleAge = 150

// SchoolConstants Holds shared Meaning, independent of Deployment.
// Keep the Snapshot private; Readers Receive a Value Copy.
type SchoolConstants struct {
	MinimumAgeYears int `json:"minimumAgeYears"`
}

var schoolConstants = settings.ReadGlobalValues[SchoolConstants](
	filepath.Join(settings.FindSchoolDataRoot(), "constants", "school.json"),
	map[string]settings.ValueRule{"minimumAgeYears": settings.CheckIntegerRange(1, OldestPlausibleAge)},
)

// ReadSchoolConstants Returns the Startup Snapshot without a Mutation Path.
func ReadSchoolConstants() SchoolConstants {
	return schoolConstants
}
