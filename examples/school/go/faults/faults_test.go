package faults

import (
	"errors"
	"fmt"
	"testing"
)

func TestEachFaultAnswersWithTheKindItDeclared(t *testing.T) {
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

func TestAnUncontrolledFailureIsOursAndItsMessageNeverLeaves(t *testing.T) {
	driver := errors.New("the driver Broke at /var/db")

	if got := ReadFaultKind(driver); got != Unexpected {
		t.Fatalf("wanted Unexpected, got %d", got)
	}
	if reason := ReadFaultReason(driver); reason != "internal error" {
		t.Fatalf("a driver Message must not Leave, got %q", reason)
	}
}

func TestAWrappedFaultStillCarriesItsKindAndReason(t *testing.T) {
	wrapped := fmt.Errorf("while Reading: %w", ReportMissingRecord("gone"))

	if ReadFaultKind(wrapped) != MissingRecord || ReadFaultReason(wrapped) != "gone" {
		t.Fatalf("kind=%d reason=%q", ReadFaultKind(wrapped), ReadFaultReason(wrapped))
	}
	if !errors.Is(wrapped, ReportMissingRecord("gone")) {
		t.Fatal("a Fault Must stay Comparable through a Wrapper")
	}
}
