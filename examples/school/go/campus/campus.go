package campus

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
// A File that is Missing, Empty, not an Array, Holds a second Value or a Text that
// is not a valid RUT Stops Startup.
func OpenOffice(path string, logger *slog.Logger) (Office, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Office{}, fmt.Errorf("registry file: %w", err)
	}

	var texts []string
	decoder := json.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&texts); err != nil || texts == nil {
		return Office{}, fmt.Errorf("registry file %s Must Hold one Array of Texts", path)
	}
	if _, err := decoder.Token(); !errors.Is(err, io.EOF) {
		return Office{}, fmt.Errorf("registry file %s Holds a second Value", path)
	}

	registered := make(map[school.RUT]bool, len(texts))
	for _, text := range texts {
		rut := school.ReadRUT(text)
		if !rut.LooksValid() {
			return Office{}, fmt.Errorf("registry file %s: %q is not a valid RUT", path, text)
		}

		registered[rut] = true
	}

	return Office{registered: registered, Logger: logger}, nil
}

// ConfirmRut Answers whether the Registry Knows a RUT. This one never Fails to Answer.
func (o Office) ConfirmRut(rut school.RUT) (bool, error) {
	return o.registered[rut], nil
}

// AnnounceEnrollment Logs the Notice, Naming the Student by Id and never by RUT or Name.
func (o Office) AnnounceEnrollment(student school.Student) error {
	o.log().Info(fmt.Sprintf("student %d is enrolled", student.ID))

	return nil
}

func (o Office) log() *slog.Logger {
	if o.Logger == nil {
		return slog.New(slog.DiscardHandler)
	}

	return o.Logger
}
