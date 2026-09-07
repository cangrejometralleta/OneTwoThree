package faults

import (
	"errors"
	"net/http"
)

// A Fault is a Failure the Program Expected, Carrying the Answer it Deserves.
// A Failure that Reaches the Edge without one is a five hundred.
type Fault struct {
	Status int
	Reason string
}

// Error Makes a Fault an ordinary Error, Comparable with errors.Is.
func (f Fault) Error() string {
	return f.Reason
}

// RefuseInvalidInput Names a Caller that Sent something the Rules Reject.
func RefuseInvalidInput(reason string) Fault {
	return Fault{Status: http.StatusBadRequest, Reason: reason}
}

// RefuseUnprovenCaller Names a Caller the Program cannot Recognise.
func RefuseUnprovenCaller(reason string) Fault {
	return Fault{Status: http.StatusUnauthorized, Reason: reason}
}

// ReportMissingRecord Names something the Caller Asked for and We Lack.
func ReportMissingRecord(reason string) Fault {
	return Fault{Status: http.StatusNotFound, Reason: reason}
}

// ReportTakenValue Names a Collision with something already Stored.
func ReportTakenValue(reason string) Fault {
	return Fault{Status: http.StatusConflict, Reason: reason}
}

// ReadFaultStatus is the one Place that Turns a Failure into a Number.
// A Fault Answers for itself; everything else Answers five hundred.
func ReadFaultStatus(err error) int {
	var fault Fault
	if errors.As(err, &fault) {
		return fault.Status
	}

	return http.StatusInternalServerError
}
