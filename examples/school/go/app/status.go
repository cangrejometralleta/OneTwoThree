package app

import (
	"net/http"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/faults"
)

// faultStatuses Translates each Kind into the Status HTTP Knows it by.
// The Core Names the Kind; only the Application Layer Knows the Protocol.
// Reference: https://www.rfc-editor.org/rfc/rfc9110#section-15
var faultStatuses = map[faults.Kind]int{
	faults.InvalidInput:        http.StatusBadRequest,
	faults.UnprovenCaller:      http.StatusUnauthorized,
	faults.MissingRecord:       http.StatusNotFound,
	faults.TakenValue:          http.StatusConflict,
	faults.UnprocessableValue:  http.StatusUnprocessableEntity,
	faults.UnavailableProvider: http.StatusServiceUnavailable,
}

// ReadFaultStatus is the one Place that Turns a Failure into a Number.
// A Kind without a Status Answers five hundred.
func ReadFaultStatus(err error) int {
	if status, known := faultStatuses[faults.ReadFaultKind(err)]; known {
		return status
	}

	return http.StatusInternalServerError
}
