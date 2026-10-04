package faults

import "errors"

// A Kind Names what Went wrong, in the Business' own Words.
// It Names no Protocol: app Translates a Kind into a Status.
type Kind int

const (
	// Unexpected is every Failure the Program did not Plan for.
	Unexpected Kind = iota
	InvalidInput
	UnprovenCaller
	MissingRecord
	TakenValue
	UnprocessableValue
	UnavailableProvider
)

// A Fault is a Failure the Program Expected, Carrying the Kind it Deserves.
// A Failure that Reaches the Edge without one is Unexpected.
type Fault struct {
	Kind   Kind
	Reason string
}

// Error Makes a Fault an ordinary Error, Comparable with errors.Is.
func (f Fault) Error() string {
	return f.Reason
}

// RefuseInvalidInput Names a Caller that Sent something the Rules Reject.
func RefuseInvalidInput(reason string) Fault {
	return Fault{Kind: InvalidInput, Reason: reason}
}

// RefuseUnprovenCaller Names a Caller the Program cannot Recognise.
func RefuseUnprovenCaller(reason string) Fault {
	return Fault{Kind: UnprovenCaller, Reason: reason}
}

// ReportMissingRecord Names something the Caller Asked for and We Lack.
func ReportMissingRecord(reason string) Fault {
	return Fault{Kind: MissingRecord, Reason: reason}
}

// ReportTakenValue Names a Collision with something already Stored.
func ReportTakenValue(reason string) Fault {
	return Fault{Kind: TakenValue, Reason: reason}
}

// RefuseUnprocessableValue Names a Value that is well Formed and still Denied
// by the Authority that Knows.
func RefuseUnprocessableValue(reason string) Fault {
	return Fault{Kind: UnprocessableValue, Reason: reason}
}

// ReportUnavailableProvider Names a Provider that cannot Answer now.
// The Caller may Retry; an Unknown is never a Yes.
func ReportUnavailableProvider(reason string) Fault {
	return Fault{Kind: UnavailableProvider, Reason: reason}
}

// ReadFaultKind Finds the Fault inside any Wrapping.
// A Fault Answers for itself; everything else is Unexpected.
func ReadFaultKind(err error) Kind {
	var fault Fault
	if errors.As(err, &fault) {
		return fault.Kind
	}

	return Unexpected
}
