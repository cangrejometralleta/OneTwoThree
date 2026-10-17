package app

import (
	"errors"
	"fmt"
	"testing"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/faults"
)

func TestReadFaultStatusAnswersForEachKind(t *testing.T) {
	cases := map[error]int{
		faults.RefuseInvalidInput("bad"):         400,
		faults.RefuseUnprovenCaller("who"):       401,
		faults.ReportMissingRecord("gone"):       404,
		faults.ReportTakenValue("taken"):         409,
		faults.RefuseUnprocessableValue("no"):    422,
		faults.ReportUnavailableProvider("down"): 503,
	}

	for fault, wanted := range cases {
		if got := ReadFaultStatus(fault); got != wanted {
			t.Errorf("%q wanted %d, got %d", fault, wanted, got)
		}
	}
}

func TestReadFaultStatusRefusesToGuess(t *testing.T) {
	if got := ReadFaultStatus(errors.New("the driver Broke")); got != 500 {
		t.Fatalf("wanted 500, got %d", got)
	}
}

func TestReadFaultStatusSurvivesWrapping(t *testing.T) {
	if got := ReadFaultStatus(fmt.Errorf("while Reading: %w", faults.ReportMissingRecord("gone"))); got != 404 {
		t.Fatalf("wanted 404 through the Wrapper, got %d", got)
	}
}
