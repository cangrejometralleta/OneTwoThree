package faults

import (
	"errors"
	"fmt"
	"testing"
)

// Every controlled Fault Answers with the Kind it Declared.
func TestReadFaultKindAnswersForEachKind(t *testing.T) {
	cases := map[error]Kind{
		RefuseInvalidInput("bad"):         InvalidInput,
		RefuseUnprovenCaller("who"):       UnprovenCaller,
		ReportMissingRecord("gone"):       MissingRecord,
		ReportTakenValue("taken"):         TakenValue,
		RefuseUnprocessableValue("no"):    UnprocessableValue,
		ReportUnavailableProvider("down"): UnavailableProvider,
	}

	for fault, wanted := range cases {
		if got := ReadFaultKind(fault); got != wanted {
			t.Errorf("%q wanted %d, got %d", fault, wanted, got)
		}
	}
}

// An uncontrolled Failure is ours.
func TestReadFaultKindRefusesToGuess(t *testing.T) {
	if got := ReadFaultKind(errors.New("the driver Broke")); got != Unexpected {
		t.Fatalf("wanted Unexpected, got %d", got)
	}
}

// A wrapped Fault still Carries its Kind.
func TestReadFaultKindSurvivesWrapping(t *testing.T) {
	wrapped := fmt.Errorf("while Reading: %w", ReportMissingRecord("gone"))

	if got := ReadFaultKind(wrapped); got != MissingRecord {
		t.Fatalf("wanted MissingRecord through the Wrapper, got %d", got)
	}

	if !errors.Is(wrapped, ReportMissingRecord("gone")) {
		t.Fatal("a Fault Must stay Comparable through a Wrapper")
	}
}
