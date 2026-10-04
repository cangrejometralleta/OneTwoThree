package app

import (
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/faults"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/wire"
)

// Telling is a Script: it Answers with a Value or Fails, and builds no Reply.
type Telling func(transport.Request) (any, error)

// AnswerWith is the one Crossing from a Telling to a Handler.
// The Route Declares the happy Status; a Fault Declares its own.
func AnswerWith(status int, tell Telling) transport.Handler {
	return func(req transport.Request) transport.Response {
		value, err := tell(req)
		if err != nil {
			return transport.Response{
				Status: ReadFaultStatus(err),
				Body:   wire.FailureView{Error: faults.ReadFaultReason(err)},
			}
		}

		return transport.Response{Status: status, Body: value}
	}
}
