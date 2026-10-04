package campus

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/school"
)

// Office Answers the Registry and Announces Enrolments: one Adapter, two Ports.
// The Registry is a JSON File of the RUTs that are Real; the Notice goes to the Log.
type Office struct {
	registered map[school.RUT]bool
	Logger     *slog.Logger
}

// OpenOffice Reads the Registry File once, at Startup.
// A File that is Missing, Broken or Holds a RUT that cannot Be one Stops Startup.
func OpenOffice(path string, logger *slog.Logger) (Office, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Office{}, fmt.Errorf("registry file: %w", err)
	}

	var ruts []school.RUT
	if err := json.Unmarshal(data, &ruts); err != nil {
		return Office{}, fmt.Errorf("registry file %s: %w", path, err)
	}

	registered := make(map[school.RUT]bool, len(ruts))
	for _, rut := range ruts {
		if !rut.LooksValid() {
			return Office{}, fmt.Errorf("registry file %s: %q is not a valid RUT", path, rut)
		}

		registered[rut] = true
	}

	return Office{registered: registered, Logger: logger}, nil
}

// ConfirmRut Answers whether the Registry Knows a RUT.
// This Registry never Fails to Answer; a Remote one would.
func (o Office) ConfirmRut(rut school.RUT) (bool, error) {
	return o.registered[rut], nil
}

// AnnounceEnrollment Sends the Notice to the Log, naming the Student by Id.
func (o Office) AnnounceEnrollment(student school.Student) error {
	o.log().Info("Enrollment Announced", "student", student.ID, "course", student.Course)

	return nil
}

func (o Office) log() *slog.Logger {
	if o.Logger == nil {
		return slog.New(slog.DiscardHandler)
	}

	return o.Logger
}
