package faults

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

// Every controlled Fault Answers with the Status it Declared.
func TestReadFaultStatusAnswersForEachKind(t *testing.T) {
	cases := map[error]int{
		RefuseInvalidInput("bad"):   http.StatusBadRequest,
		RefuseUnprovenCaller("who"): http.StatusUnauthorized,
		ReportMissingRecord("gone"): http.StatusNotFound,
		ReportTakenValue("taken"):   http.StatusConflict,
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
	wrapped := fmt.Errorf("while Reading: %w", ReportMissingRecord("gone"))

	if got := ReadFaultStatus(wrapped); got != http.StatusNotFound {
		t.Fatalf("wanted 404 through the Wrapper, got %d", got)
	}

	if !errors.Is(wrapped, ReportMissingRecord("gone")) {
		t.Fatal("a Fault Must stay Comparable through a Wrapper")
	}
}
