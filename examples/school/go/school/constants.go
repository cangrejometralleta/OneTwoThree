package school

import (
	"path/filepath"
	"sync"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/settings"
)

// SchoolConstants Holds the Global Values no Deployment can Override.
type SchoolConstants struct {
	MinimumAgeYears int `json:"minimumAgeYears"`
}

var readConstants = sync.OnceValue(func() SchoolConstants {
	path := filepath.Join(settings.FindSchoolDataRoot(), "constants", "school.json")
	rules := map[string]settings.ValueRule{"minimumAgeYears": settings.CheckIntegerRange(1, 200)}

	return settings.ReadGlobalValues[SchoolConstants](path, rules)
})

// ReadSchoolConstants Loads the Constants once; a Broken File Stops the Program.
func ReadSchoolConstants() SchoolConstants {
	return readConstants()
}
