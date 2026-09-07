package app

import (
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/faults"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
)

// A Telling Speaks Business only: it Answers with a Value, or it Fails.
// It Names no Status and Builds no Reply.
type Telling func(transport.Request) (any, error)

// AnswerWith Turns a Telling into a Handler the Transport can Mount.
// The Route Declares the happy Status; a Fault Declares its own.
// This is the only Place in the Program that Builds a Reply.
func AnswerWith(status int, tell Telling) transport.Handler {
	return func(req transport.Request) transport.Response {
		body, err := tell(req)
		if err != nil {
			return transport.Response{
				Status: faults.ReadFaultStatus(err),
				Body:   map[string]string{"error": err.Error()},
			}
		}

		return transport.Response{Status: status, Body: body}
	}
}
