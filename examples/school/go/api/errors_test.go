package api

import (
	"errors"
	"net/http"
	"testing"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/app"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/faults"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/school"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/tokens"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
)

// Every Fault the Service Declares, Reached the Way a Caller Reaches it.
// Each Name is a Use Case, because a controlled Failure is one.
func TestEachFaultReachesTheEdgeWhole(t *testing.T) {
	cases := []struct {
		story  string
		arrive func(SchoolAPI) error
		want   error
		status int
	}{
		{
			story:  "someone Enrols with no Name",
			arrive: EnrolWithBody(`{"rut":"12345678-5","name":"","age":20,"courseId":1}`),
			want:   school.ErrNameIsEmpty,
			status: http.StatusBadRequest,
		},
		{
			story:  "someone Enrols with a RUT that Fails its Digit",
			arrive: EnrolWithBody(`{"rut":"12345678-9","name":"Ada","age":20,"courseId":1}`),
			want:   school.ErrRutIsInvalid,
			status: http.StatusBadRequest,
		},
		{
			story:  "someone Enrols below the Enrolment Age",
			arrive: EnrolWithBody(`{"rut":"12345678-5","name":"Ada","age":9,"courseId":1}`),
			want:   school.ErrAgeIsTooLow,
			status: http.StatusBadRequest,
		},
		{
			story:  "someone Enrols into a Course that never Opened",
			arrive: EnrolWithBody(`{"rut":"12345678-5","name":"Ada","age":20,"courseId":99}`),
			want:   school.ErrCourseUnknown,
			status: http.StatusNotFound,
		},
		{
			story:  "someone Sends a Body that is not JSON",
			arrive: EnrolWithBody(`{"rut":`),
			want:   app.ErrBodyIsBroken,
			status: http.StatusBadRequest,
		},
		{
			story: "someone Enrols with a RUT already Registered",
			arrive: func(api SchoolAPI) error {
				EnrolWithBody(`{"rut":"12345678-5","name":"Ada","age":20,"courseId":1}`)(api)

				return EnrolWithBody(`{"rut":"12345678-5","name":"Grace","age":22,"courseId":1}`)(api)
			},
			want:   school.ErrRutTaken,
			status: http.StatusConflict,
		},
		{
			story: "someone Asks for a Student by a Path that Holds no Number",
			arrive: func(api SchoolAPI) error {
				return TellingFailure(api.ShowStudentRecord, transport.Request{Path: map[string]string{"id": "abc"}})
			},
			want:   app.ErrPathIsBroken,
			status: http.StatusBadRequest,
		},
		{
			story: "someone Asks for a Student who never Enrolled",
			arrive: func(api SchoolAPI) error {
				return TellingFailure(api.ShowStudentRecord, transport.Request{Path: map[string]string{"id": "42"}})
			},
			want:   school.ErrStudentUnknown,
			status: http.StatusNotFound,
		},
		{
			story: "someone Asks for a Course that never Opened",
			arrive: func(api SchoolAPI) error {
				return TellingFailure(api.ShowCourseRecord, transport.Request{Path: map[string]string{"id": "99"}})
			},
			want:   school.ErrCourseUnknown,
			status: http.StatusNotFound,
		},
		{
			story: "someone Asks for a Page that Counts backwards",
			arrive: func(api SchoolAPI) error {
				return TellingFailure(api.ListStudentRecords, transport.Request{Query: map[string]string{"page": "-1"}})
			},
			want:   school.ErrPageIsInvalid,
			status: http.StatusBadRequest,
		},
		{
			story: "someone Drops a Student who already Left",
			arrive: func(api SchoolAPI) error {
				return TellingFailure(api.DropStudentRecord, transport.Request{Path: map[string]string{"id": "7"}})
			},
			want:   school.ErrStudentUnknown,
			status: http.StatusNotFound,
		},
		{
			story: "someone Rewrites an Enrolment that never Existed",
			arrive: func(api SchoolAPI) error {
				return TellingFailure(api.SaveStudentRecord, transport.Request{
					Path: map[string]string{"id": "7"},
					Body: []byte(`{"rut":"12345678-5","name":"Ada","age":20,"courseId":1}`),
				})
			},
			want:   school.ErrStudentUnknown,
			status: http.StatusNotFound,
		},
		{
			story: "someone Arrives with no Token at all",
			arrive: func(api SchoolAPI) error {
				guarded := app.RequireProvenCaller(api.Tokens, api.ListStudentRecords)

				return TellingFailure(guarded, transport.Request{})
			},
			want:   tokens.ErrTokenIsInvalid,
			status: http.StatusUnauthorized,
		},
	}

	for _, test := range cases {
		t.Run(test.story, func(t *testing.T) {
			CheckFault(t, test.arrive(BuildTestingSchool()), test.want, test.status)
		})
	}
}

// EnrolWithBody Spells the most common Arrival once.
func EnrolWithBody(body string) func(SchoolAPI) error {
	return func(api SchoolAPI) error {
		return TellingFailure(api.AddStudentRecord, transport.Request{Body: []byte(body)})
	}
}

// TellingFailure Keeps only the Half a Fault Test Cares about.
func TellingFailure(tell app.Telling, req transport.Request) error {
	_, err := tell(req)

	return err
}

// CheckFault Reads both Halves: the Fault that Arrived and the Answer it Carries.
// The Number is Spelled in the Table, never Read from the Fault under Test,
// so a Declaration that Drifts Fails here instead of Agreeing with itself.
func CheckFault(t *testing.T, got, want error, status int) {
	t.Helper()

	if !errors.Is(got, want) {
		t.Fatalf("wanted %v, got %v", want, got)
	}

	if answer := faults.ReadFaultStatus(got); answer != status {
		t.Errorf("wanted %d, got %d", status, answer)
	}
}
