package app

import (
	"errors"
	"testing"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/faults"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/wire"
)

type fixedReader struct{ caller string }

func (r fixedReader) ReadCaller(token string) (string, error) {
	if token != "good" {
		return "", faults.RefuseUnprovenCaller("token is Invalid or Expired")
	}

	return r.caller, nil
}

func TestRequireProvenCallerNamesTheCallerBeforeTheStoryStarts(t *testing.T) {
	var seen string
	guarded := RequireProvenCaller(fixedReader{caller: "student-registry"}, func(req transport.Request) (any, error) {
		seen = req.Caller

		return "told", nil
	})

	if _, err := guarded(transport.Request{Token: "good"}); err != nil || seen != "student-registry" {
		t.Fatalf("seen=%q err=%v", seen, err)
	}
}

func TestRequireProvenCallerNeverStartsTheStoryForAStranger(t *testing.T) {
	started := false
	guarded := RequireProvenCaller(fixedReader{}, func(transport.Request) (any, error) {
		started = true

		return nil, nil
	})

	if _, err := guarded(transport.Request{Token: "forged"}); ReadFaultStatus(err) != 401 || started {
		t.Fatalf("err=%v started=%v", err, started)
	}
}

func TestAnswerWithCarriesTheDeclaredStatusOrTheFaultOwn(t *testing.T) {
	happy := AnswerWith(201, func(transport.Request) (any, error) { return "made", nil })
	if reply := happy(transport.Request{}); reply.Status != 201 || reply.Body != "made" {
		t.Fatalf("reply=%v", reply)
	}

	failing := AnswerWith(200, func(transport.Request) (any, error) {
		return nil, faults.ReportMissingRecord("gone")
	})
	reply := failing(transport.Request{})
	if reply.Status != 404 || reply.Body != (wire.FailureView{Error: "gone"}) {
		t.Fatalf("reply=%v", reply)
	}

	ours := AnswerWith(200, func(transport.Request) (any, error) { return nil, errors.New("driver at /var/db broke") })
	reply = ours(transport.Request{})
	if reply.Status != 500 || reply.Body != (wire.FailureView{Error: "internal error"}) {
		t.Fatalf("a driver Message must not Leave, reply=%v", reply)
	}
}
