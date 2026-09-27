package app

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/faults"
)

// Every controlled Fault Answers with the Status its Kind Maps to.
func TestReadFaultStatusAnswersForEachKind(t *testing.T) {
	cases := map[error]int{
		faults.RefuseInvalidInput("bad"):   http.StatusBadRequest,
		faults.RefuseUnprovenCaller("who"): http.StatusUnauthorized,
		faults.ReportMissingRecord("gone"): http.StatusNotFound,
		faults.ReportTakenValue("taken"):   http.StatusConflict,
	}

	for fault, wanted := range cases {
		if got := ReadFaultStatus(fault); got != wanted {
			t.Errorf("%q wanted %d, got %d", fault, wanted, got)
		}
	}
}

// An uncontrolled Failure is ours.
func TestReadFaultStatusRefusesToGuess(t *testing.T) {
	if got := ReadFaultStatus(errors.New("the driver Broke")); got != http.StatusInternalServerError {
		t.Fatalf("wanted 500, got %d", got)
	}
}

// A wrapped Fault still Carries its Answer.
func TestReadFaultStatusSurvivesWrapping(t *testing.T) {
	wrapped := fmt.Errorf("while Reading: %w", faults.ReportMissingRecord("gone"))

	if got := ReadFaultStatus(wrapped); got != http.StatusNotFound {
		t.Fatalf("wanted 404 through the Wrapper, got %d", got)
	}
}
