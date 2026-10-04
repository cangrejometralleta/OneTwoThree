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

// An uncontrolled Failure is ours, and its Message never Leaves.
func TestReadFaultKindRefusesToGuess(t *testing.T) {
	driver := errors.New("the driver Broke at /var/db")

	if got := ReadFaultKind(driver); got != Unexpected {
		t.Fatalf("wanted Unexpected, got %d", got)
	}
	if reason := ReadFaultReason(driver); reason != "internal error" {
		t.Fatalf("a driver Message must not Leave, got %q", reason)
	}
}

// A wrapped Fault still Carries its Kind and its Reason.
func TestReadFaultKindSurvivesWrapping(t *testing.T) {
	wrapped := fmt.Errorf("while Reading: %w", ReportMissingRecord("gone"))

	if got := ReadFaultKind(wrapped); got != MissingRecord {
		t.Fatalf("wanted MissingRecord through the Wrapper, got %d", got)
	}
	if reason := ReadFaultReason(wrapped); reason != "gone" {
		t.Fatalf("wanted the Declared Reason, got %q", reason)
	}
	if !errors.Is(wrapped, ReportMissingRecord("gone")) {
		t.Fatal("a Fault Must stay Comparable through a Wrapper")
	}
}
