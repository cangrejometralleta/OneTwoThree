package app

import (
	"errors"
	"testing"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/faults"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/wire"
)

type fixedReader struct{}

func (fixedReader) ReadCaller(token string) (string, error) {
	if token != "good" {
		return "", faults.RefuseUnprovenCaller("token is Invalid or Expired")
	}

	return "student-registry", nil
}

func bearing(header string) transport.Request {
	return transport.Request{Headers: map[string]string{"authorization": header}}
}

func TestTheCallerIsNamedBeforeTheStoryStartsWhateverTheCaseOfTheScheme(t *testing.T) {
	var seen string
	guarded := RequireProvenCaller(fixedReader{}, func(req transport.Request) (any, error) {
		seen = req.Caller

		return "told", nil
	})

	for _, header := range []string{"Bearer good", "bearer good", "BEARER good"} {
		seen = ""
		if _, err := guarded(bearing(header)); err != nil || seen != "student-registry" {
			t.Errorf("%s: seen=%q err=%v", header, seen, err)
		}
	}
}

func TestTheStoryNeverStartsForAStranger(t *testing.T) {
	started := false
	guarded := RequireProvenCaller(fixedReader{}, func(transport.Request) (any, error) {
		started = true

		return nil, nil
	})

	if reply := AnswerWith(200, guarded)(bearing("Bearer forged")); reply.Status != 401 || started {
		t.Fatalf("reply=%v started=%v", reply, started)
	}
}

func TestAnAnswerCarriesTheDeclaredStatusOrTheFaultsOwn(t *testing.T) {
	if reply := AnswerWith(201, func(transport.Request) (any, error) { return "made", nil })(transport.Request{}); reply != (transport.Response{Status: 201, Body: "made"}) {
		t.Fatalf("reply=%v", reply)
	}

	failing := AnswerWith(200, func(transport.Request) (any, error) { return nil, faults.ReportMissingRecord("gone") })
	if reply := failing(transport.Request{}); reply != (transport.Response{Status: 404, Body: wire.FailureView{Error: "gone"}}) {
		t.Fatalf("reply=%v", reply)
	}

	ours := AnswerWith(200, func(transport.Request) (any, error) { return nil, errors.New("driver at /var/db broke") })
	if reply := ours(transport.Request{}); reply != (transport.Response{Status: 500, Body: wire.FailureView{Error: "internal error"}}) {
		t.Fatalf("a driver Message must not Leave, reply=%v", reply)
	}
}
