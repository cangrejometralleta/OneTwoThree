package main

// SchoolConstants Holds shared Meaning, independent of Deployment.
// Keep the Snapshot private; Readers Receive a Value Copy.
type SchoolConstants struct {
	MinimumAgeYears int `json:"minimumAgeYears"`
}

var schoolConstants = readGlobalValues[SchoolConstants]("../constants/school.json", map[string]valueRule{
	"minimumAgeYears": checkIntegerRange(1, 150),
})

// ReadSchoolConstants Returns the Startup Snapshot without a Mutation Path.
func ReadSchoolConstants() SchoolConstants {
	return schoolConstants
}
