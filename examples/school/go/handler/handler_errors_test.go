package handler

import (
	"testing"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/wire"
)

// Every Fault the Spec Names, Reached the Way a Caller Reaches it.
// The Status and the Reason are Spelled here, never Read from the Fault under Test,
// so a Declaration that Drifts Fails here instead of Agreeing with itself.
func TestEachFaultReachesTheEdgeWhole(t *testing.T) {
	enrolWith := func(body string) func(*testing.T, Handler, *FakeOffice) transport.Response {
		return func(t *testing.T, api Handler, _ *FakeOffice) transport.Response {
			return CallRoute(t, api, "POST", "/students", transport.Request{Body: []byte(body)})
		}
	}
	showing := func(pattern, id string) func(*testing.T, Handler, *FakeOffice) transport.Response {
		return func(t *testing.T, api Handler, _ *FakeOffice) transport.Response {
			return CallRoute(t, api, "GET", pattern, transport.Request{Path: map[string]string{"id": id}})
		}
	}

	cases := []struct {
		story  string
		arrive func(*testing.T, Handler, *FakeOffice) transport.Response
		status int
		reason string
	}{
		{"someone Enrols with no Name", enrolWith(`{"rut":"12345678-5","name":"","age":20,"courseId":1}`), 400, "name is Empty"},
		{"someone Opens a Course with no Code", func(t *testing.T, api Handler, _ *FakeOffice) transport.Response {
			return CallRoute(t, api, "POST", "/courses", transport.Request{Body: []byte(`{"code":" ","name":"Physics"}`)})
		}, 400, "code is Empty"},
		{"someone Enrols with a RUT that Fails its Digit", enrolWith(`{"rut":"12345678-9","name":"Ada","age":20,"courseId":1}`), 400, "rut Fails its Check Digit"},
		{"someone Enrols below the Enrolment Age", enrolWith(`{"rut":"12345678-5","name":"Ada","age":9,"courseId":1}`), 400, "age Must be 18 or more"},
		{"someone Asks for a Page that Counts backwards", func(t *testing.T, api Handler, _ *FakeOffice) transport.Response {
			return CallRoute(t, api, "GET", "/students", transport.Request{Query: map[string]string{"page": "-1"}})
		}, 400, "page or size is not valid"},
		{"someone Sends a Body that is not JSON", enrolWith(`{"rut":`), 400, "body is not valid JSON"},
		{"someone Sends a Body with an Identity of their own", enrolWith(`{"id":5,"rut":"12345678-5","name":"Ada","age":20,"courseId":1}`), 400, "body is not valid JSON"},
		{"someone Asks for a Student by a Path that Holds no Number", showing("/students/{id}", "abc"), 400, "path Holds no valid Identity"},
		{"someone Arrives with no Token at all", func(t *testing.T, api Handler, _ *FakeOffice) transport.Response {
			return api.DeclareSchoolRoutes().Routes[1].Handle(transport.Request{})
		}, 401, "token is Invalid or Expired"},
		{"someone Asks for a Student who never Enrolled", showing("/students/{id}", "42"), 404, "student not Found"},
		{"someone Asks for a Student beyond 64 bits", showing("/students/{id}", "99999999999999999999"), 404, "student not Found"},
		{"someone Asks for a Course that never Opened", showing("/courses/{id}", "99"), 404, "course not Found"},
		{"someone Enrols into a Course that never Opened", enrolWith(`{"rut":"12345678-5","name":"Ada","age":20,"courseId":-1}`), 404, "course not Found"},
		{"someone Knocks on a Door the Libretto does not Hold", func(t *testing.T, api Handler, _ *FakeOffice) transport.Response {
			return api.DeclareSchoolRoutes().Unknown(transport.Request{})
		}, 404, "route not Found"},
		{"someone Enrols with a RUT already Registered", func(t *testing.T, api Handler, front *FakeOffice) transport.Response {
			enrolWith(`{"rut":"12345678-5","name":"Ada","age":20,"courseId":1}`)(t, api, front)

			return enrolWith(`{"rut":"12345678-5","name":"Grace","age":22,"courseId":1}`)(t, api, front)
		}, 409, "rut is already Registered"},
		{"someone Enrols with a RUT the Registry Denies", func(t *testing.T, api Handler, front *FakeOffice) transport.Response {
			front.Denies = true

			return enrolWith(enrolBody)(t, api, front)
		}, 422, "rut is not in the Registry"},
		{"someone Rewrites a Student nobody Holds with a RUT the Registry Denies", func(t *testing.T, api Handler, front *FakeOffice) transport.Response {
			front.Denies = true

			return CallRoute(t, api, "PUT", "/students/{id}", transport.Request{Path: map[string]string{"id": "7"}, Body: []byte(enrolBody)})
		}, 422, "rut is not in the Registry"},
		{"someone Enrols while the Registry cannot Answer", func(t *testing.T, api Handler, front *FakeOffice) transport.Response {
			front.Silent = true

			return enrolWith(enrolBody)(t, api, front)
		}, 503, "registry cannot Answer now"},
		{"someone Drops a Student who already Left", func(t *testing.T, api Handler, _ *FakeOffice) transport.Response {
			return CallRoute(t, api, "DELETE", "/students/{id}", transport.Request{Path: map[string]string{"id": "7"}})
		}, 404, "student not Found"},
		{"someone Rewrites an Enrolment that never Existed", func(t *testing.T, api Handler, _ *FakeOffice) transport.Response {
			return CallRoute(t, api, "PUT", "/students/{id}", transport.Request{Path: map[string]string{"id": "7"}, Body: []byte(enrolBody)})
		}, 404, "student not Found"},
	}

	for _, test := range cases {
		t.Run(test.story, func(t *testing.T) {
			api, front := BuildTestingSchool()
			reply := test.arrive(t, api, front)
			if reply != (transport.Response{Status: test.status, Body: wire.FailureView{Error: test.reason}}) {
				t.Fatalf("wanted %d %q, got %d %v", test.status, test.reason, reply.Status, reply.Body)
			}
		})
	}
}
