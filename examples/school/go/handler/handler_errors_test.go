package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/app"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/school"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/tokens"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
)

const validBody = `{"rut":"12345678-5","name":"Ada","age":20,"courseId":1}`

// Every Fault the Spec Names, Reached the Way a Caller Reaches it.
// Each Name is a Use Case, because a controlled Failure is one.
func TestEachFaultReachesTheEdgeWhole(t *testing.T) {
	cases := []struct {
		story  string
		arrive func(Handler) error
		want   error
		status int
	}{
		{"someone Enrols with no Name", enrolWith(`{"rut":"12345678-5","name":"","age":20,"courseId":1}`), school.ErrNameIsEmpty, http.StatusBadRequest},
		{"someone Enrols with a RUT that Fails its Digit", enrolWith(`{"rut":"12345678-9","name":"Ada","age":20,"courseId":1}`), school.ErrRutIsInvalid, http.StatusBadRequest},
		{"someone Enrols below the Enrolment Age", enrolWith(`{"rut":"12345678-5","name":"Ada","age":9,"courseId":1}`), school.ErrAgeIsTooLow, http.StatusBadRequest},
		{"someone Enrols into a Course that never Opened", enrolWith(`{"rut":"12345678-5","name":"Ada","age":20,"courseId":99}`), school.ErrCourseUnknown, http.StatusNotFound},
		{"someone Sends a Body that is not JSON", enrolWith(`{"rut":`), app.ErrBodyIsBroken, http.StatusBadRequest},
		{"someone Sends a Body with a Field Missing", enrolWith(`{"rut":"12345678-5","name":"Ada","age":20}`), app.ErrBodyIsBroken, http.StatusBadRequest},
		{"someone Sends a Body with an Identity of their own", enrolWith(`{"id":5,"rut":"12345678-5","name":"Ada","age":20,"courseId":1}`), app.ErrBodyIsBroken, http.StatusBadRequest},
		{
			"someone Enrols with a RUT already Registered",
			func(api Handler) error {
				enrolWith(validBody)(api)

				return enrolWith(`{"rut":"12345678-5","name":"Grace","age":22,"courseId":1}`)(api)
			},
			school.ErrRutTaken, http.StatusConflict,
		},
		{
			"someone Enrols with a RUT the Registry Denies",
			func(api Handler) error {
				api.School.Registry = DenyingOffice{}

				return enrolWith(validBody)(api)
			},
			school.ErrRutUnregistered, http.StatusUnprocessableEntity,
		},
		{
			"someone Enrols while the Registry cannot Answer",
			func(api Handler) error {
				api.School.Registry = SilentOffice{}

				return enrolWith(validBody)(api)
			},
			school.ErrRegistryUnavailable, http.StatusServiceUnavailable,
		},
		{"someone Opens a Course with no Code", telling(Handler.AddCourseRecord, transport.Request{Body: []byte(`{"code":"","name":"Physics"}`)}), school.ErrCodeIsEmpty, http.StatusBadRequest},
		{"someone Opens a Course with no Name", telling(Handler.AddCourseRecord, transport.Request{Body: []byte(`{"code":"FIS-201","name":""}`)}), school.ErrNameIsEmpty, http.StatusBadRequest},
		{"someone Asks for a Student by a Path that Holds no Number", telling(Handler.ShowStudentRecord, pathID("abc")), app.ErrPathIsBroken, http.StatusBadRequest},
		{"someone Asks for a Student by the Number zero", telling(Handler.ShowStudentRecord, pathID("0")), app.ErrPathIsBroken, http.StatusBadRequest},
		{"someone Asks for a Student who never Enrolled", telling(Handler.ShowStudentRecord, pathID("42")), school.ErrStudentUnknown, http.StatusNotFound},
		{"someone Asks for a Course that never Opened", telling(Handler.ShowCourseRecord, pathID("99")), school.ErrCourseUnknown, http.StatusNotFound},
		{"someone Asks for a Page that Counts backwards", telling(Handler.ListStudentRecords, queryOf("page", "-1")), school.ErrPageIsInvalid, http.StatusBadRequest},
		{"someone Asks for a Page Numbered with a Word", telling(Handler.ListStudentRecords, queryOf("page", "abc")), school.ErrPageIsInvalid, http.StatusBadRequest},
		{"someone Asks for a Size Measured in Words", telling(Handler.ListCourseRecords, queryOf("size", "many")), school.ErrPageIsInvalid, http.StatusBadRequest},
		{"someone Asks for a Size of zero", telling(Handler.ListCourseRecords, queryOf("size", "0")), school.ErrPageIsInvalid, http.StatusBadRequest},
		{"someone Drops a Student who already Left", telling(Handler.DropStudentRecord, pathID("7")), school.ErrStudentUnknown, http.StatusNotFound},
		{
			"someone Rewrites an Enrolment that never Existed",
			telling(Handler.SaveStudentRecord, transport.Request{Path: map[string]string{"id": "7"}, Body: []byte(validBody)}),
			school.ErrStudentUnknown, http.StatusNotFound,
		},
		{
			"someone Arrives with no Token at all",
			func(api Handler) error {
				guarded := app.RequireProvenCaller(api.Tokens, api.ListStudentRecords)

				return tellingFailure(guarded, transport.Request{})
			},
			tokens.ErrTokenIsInvalid, http.StatusUnauthorized,
		},
	}

	for _, test := range cases {
		t.Run(test.story, func(t *testing.T) {
			CheckFault(t, test.arrive(BuildTestingSchool()), test.want, test.status)
		})
	}
}

func enrolWith(body string) func(Handler) error {
	return telling(Handler.AddStudentRecord, transport.Request{Body: []byte(body)})
}

func telling(script func(Handler, transport.Request) (any, error), req transport.Request) func(Handler) error {
	return func(api Handler) error {
		_, err := script(api, req)

		return err
	}
}

func pathID(id string) transport.Request { return transport.Request{Path: map[string]string{"id": id}} }

func queryOf(name, value string) transport.Request {
	return transport.Request{Query: map[string]string{name: value}}
}

func tellingFailure(tell app.Telling, req transport.Request) error {
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

	if answer := app.ReadFaultStatus(got); answer != status {
		t.Errorf("wanted %d, got %d", status, answer)
	}
}
