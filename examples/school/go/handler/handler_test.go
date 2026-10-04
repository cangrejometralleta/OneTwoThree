package handler

import (
	"errors"
	"net/http"
	"sort"
	"testing"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/app"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/school"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/tokens"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/wire"
)

// FakeSchool Stands in for the Store. No Database, no Framework, no Port:
// the Providers Allow it. Course 1 is Open, and the RUT is Unique like a real Index.
type FakeSchool struct {
	students map[school.StudentID]school.Student
	courses  map[school.CourseID]school.Course
}

func BuildFakeSchool() *FakeSchool {
	return &FakeSchool{
		students: map[school.StudentID]school.Student{},
		courses:  map[school.CourseID]school.Course{1: {ID: 1, Code: "MAT-101", Name: "Algebra"}},
	}
}

func (f *FakeSchool) InsertStudentRow(s school.Student) (school.Student, error) {
	for _, enrolled := range f.students {
		if enrolled.Rut == s.Rut {
			return school.Student{}, school.ErrRutTaken
		}
	}
	s.ID = school.StudentID(len(f.students) + 1)
	f.students[s.ID] = s

	return s, nil
}

func (f *FakeSchool) SelectStudentRow(id school.StudentID) (school.Student, error) {
	student, found := f.students[id]
	if !found {
		return school.Student{}, school.ErrStudentUnknown
	}

	return student, nil
}

func (f *FakeSchool) UpdateStudentRow(s school.Student) (school.Student, error) {
	if _, found := f.students[s.ID]; !found {
		return school.Student{}, school.ErrStudentUnknown
	}
	f.students[s.ID] = s

	return s, nil
}

func (f *FakeSchool) DeleteStudentRow(id school.StudentID) error {
	if _, found := f.students[id]; !found {
		return school.ErrStudentUnknown
	}
	delete(f.students, id)

	return nil
}

func (f *FakeSchool) SelectStudentPage(school.Page) ([]school.Student, error) {
	roll := make([]school.Student, 0, len(f.students))
	for _, student := range f.students {
		roll = append(roll, student)
	}
	sort.Slice(roll, func(a, b int) bool { return roll[a].ID < roll[b].ID })

	return roll, nil
}

func (f *FakeSchool) InsertCourseRow(c school.Course) (school.Course, error) {
	c.ID = school.CourseID(len(f.courses) + 1)
	f.courses[c.ID] = c

	return c, nil
}

func (f *FakeSchool) SelectCourseRow(id school.CourseID) (school.Course, error) {
	course, found := f.courses[id]
	if !found {
		return school.Course{}, school.ErrCourseUnknown
	}

	return course, nil
}

func (f *FakeSchool) SelectCoursePage(school.Page) ([]school.Course, error) {
	catalogue := make([]school.Course, 0, len(f.courses))
	for _, course := range f.courses {
		catalogue = append(catalogue, course)
	}
	sort.Slice(catalogue, func(a, b int) bool { return catalogue[a].ID < catalogue[b].ID })

	return catalogue, nil
}

// FakeOffice Stands in for the Registry and the Notice: every RUT is Real.
type FakeOffice struct{}

func (FakeOffice) ConfirmRut(school.RUT) (bool, error) { return true, nil }

func (FakeOffice) AnnounceEnrollment(school.Student) error { return nil }

// DenyingOffice Stands in for a Registry that Knows no one.
type DenyingOffice struct{ FakeOffice }

func (DenyingOffice) ConfirmRut(school.RUT) (bool, error) { return false, nil }

// SilentOffice Stands in for a Registry that cannot Answer, and a Notice that cannot Leave.
type SilentOffice struct{ FakeOffice }

func (SilentOffice) ConfirmRut(school.RUT) (bool, error) {
	return false, errors.New("registry socket closed")
}

func (SilentOffice) AnnounceEnrollment(school.Student) error {
	return errors.New("mail is down")
}

// BuildTestingSchool Hands the API its Fakes.
func BuildTestingSchool() Handler {
	fake := BuildFakeSchool()

	return Handler{
		School: app.SchoolService{Students: fake, Courses: fake, Registry: FakeOffice{}, Notices: FakeOffice{}},
		Tokens: tokens.BuildAccessTokens("test", 60),
	}
}

// CallRoute Arrives the Way a Client Arrives: through the Libretto.
// It Crosses AnswerWith and the Guard, so a Test Sees the real Status.
func CallRoute(t *testing.T, api Handler, method, pattern string, req transport.Request) transport.Response {
	t.Helper()

	token, err := api.Tokens.IssueAccessToken("student-registry")
	if err != nil {
		t.Fatal(err)
	}
	req.Token = token

	for _, route := range api.DeclareSchoolRoutes() {
		if route.Method == method && route.Pattern == pattern {
			return route.Handle(req)
		}
	}

	t.Fatalf("no Route Answers %s %s", method, pattern)

	return transport.Response{}
}

// The happy Status Lives in the Route.
func TestEachRouteAnswersWithItsDeclaredStatus(t *testing.T) {
	api := BuildTestingSchool()
	enrol := transport.Request{Body: []byte(`{"rut":"12345678-5","name":"Ada","age":20,"courseId":1}`)}
	rewrite := transport.Request{
		Path: map[string]string{"id": "1"},
		Body: []byte(`{"rut":"12345678-5","name":"Ada Lovelace","age":21,"courseId":1}`),
	}
	open := transport.Request{Body: []byte(`{"code":"FIS-201","name":"Physics"}`)}
	one := transport.Request{Path: map[string]string{"id": "1"}}

	cases := []struct {
		story   string
		method  string
		pattern string
		req     transport.Request
		status  int
	}{
		{"a Token is Minted", "POST", "/token", transport.Request{}, http.StatusCreated},
		{"someone Enrols", "POST", "/students", enrol, http.StatusCreated},
		{"the Roll is Read", "GET", "/students", transport.Request{}, http.StatusOK},
		{"one Student is Read", "GET", "/students/{id}", one, http.StatusOK},
		{"an Enrolment is Rewritten", "PUT", "/students/{id}", rewrite, http.StatusOK},
		{"the Catalogue is Read", "GET", "/courses", transport.Request{}, http.StatusOK},
		{"a Course Opens", "POST", "/courses", open, http.StatusCreated},
		{"one Course is Read", "GET", "/courses/{id}", one, http.StatusOK},
		{"an Enrolment Ends", "DELETE", "/students/{id}", one, http.StatusNoContent},
	}

	for _, test := range cases {
		t.Run(test.story, func(t *testing.T) {
			if reply := CallRoute(t, api, test.method, test.pattern, test.req); reply.Status != test.status {
				t.Fatalf("wanted %d, got %d: %v", test.status, reply.Status, reply.Body)
			}
		})
	}
}

// Every Route the Libretto Declares has a Story above, and the Spec Names nine.
func TestTheLibrettoDeclaresExactlyTheNineRoutesOfTheSpec(t *testing.T) {
	routes := BuildTestingSchool().DeclareSchoolRoutes()

	if len(routes) != 9 {
		t.Fatalf("the Spec Names nine Routes, the Libretto Declares %d", len(routes))
	}
}

func TestARewriteChangesTheStudentAndADeleteLeavesNoBody(t *testing.T) {
	api := BuildTestingSchool()
	CallRoute(t, api, "POST", "/students", transport.Request{Body: []byte(`{"rut":"12345678-5","name":"Ada","age":20,"courseId":1}`)})

	rewritten := CallRoute(t, api, "PUT", "/students/{id}", transport.Request{
		Path: map[string]string{"id": "1"},
		Body: []byte(`{"rut":"12345678-5","name":"Ada Lovelace","age":21,"courseId":1}`),
	})
	view, ok := rewritten.Body.(wire.StudentView)
	if !ok || view.Name != "Ada Lovelace" || view.Age != 21 {
		t.Fatalf("the Rewrite must Answer the Student as it Stands, got %v", rewritten.Body)
	}

	dropped := CallRoute(t, api, "DELETE", "/students/{id}", transport.Request{Path: map[string]string{"id": "1"}})
	if dropped.Status != http.StatusNoContent || dropped.Body != nil {
		t.Fatalf("a Delete Answers 204 and no Body, got %d %v", dropped.Status, dropped.Body)
	}
}

func TestAFailureReachesTheClientAsAFailureViewWithItsReason(t *testing.T) {
	api := BuildTestingSchool()

	reply := CallRoute(t, api, "GET", "/students/{id}", transport.Request{Path: map[string]string{"id": "42"}})

	if reply.Status != http.StatusNotFound || reply.Body != (wire.FailureView{Error: "student not Found"}) {
		t.Fatalf("reply=%v", reply)
	}
}

// A Notice that did not Leave never Fails the Enrolment the Caller Asked for.
func TestEnrolmentSurvivesAFailedAnnouncement(t *testing.T) {
	api := BuildTestingSchool()
	api.School.Notices = SilentOffice{}
	enrol := transport.Request{Body: []byte(`{"rut":"12345678-5","name":"Ada","age":20,"courseId":1}`)}

	if reply := CallRoute(t, api, "POST", "/students", enrol); reply.Status != http.StatusCreated {
		t.Fatalf("wanted %d, got %d: %v", http.StatusCreated, reply.Status, reply.Body)
	}
}

// Every Route but the Token Names its Caller first.
func TestEveryRouteButTheTokenRefusesAnUnnamedCaller(t *testing.T) {
	api := BuildTestingSchool()

	for _, route := range api.DeclareSchoolRoutes() {
		reply := route.Handle(transport.Request{Path: map[string]string{"id": "1"}})
		if route.Pattern == "/token" {
			continue
		}
		if reply.Status != http.StatusUnauthorized || reply.Body != (wire.FailureView{Error: "token is Invalid or Expired"}) {
			t.Errorf("%s %s must Refuse an unnamed caller, got %d %v", route.Method, route.Pattern, reply.Status, reply.Body)
		}
	}
}

func TestGuardedHandlerReceivesTheNamedCaller(t *testing.T) {
	api := BuildTestingSchool()
	var seen string
	guarded := api.Guarded(http.StatusOK, func(req transport.Request) (any, error) {
		seen = req.Caller

		return "told", nil
	})
	token, _ := api.Tokens.IssueAccessToken("student-registry")

	if reply := guarded(transport.Request{Token: token}); reply.Status != http.StatusOK || seen != "student-registry" {
		t.Fatalf("reply=%v seen=%q", reply, seen)
	}
}

func TestUncontrolledFailureAnswersFiveHundred(t *testing.T) {
	api := BuildTestingSchool()
	broken := api.Guarded(http.StatusOK, func(transport.Request) (any, error) {
		return nil, errors.New("the driver Broke at /var/db")
	})
	token, _ := api.Tokens.IssueAccessToken("student-registry")

	reply := broken(transport.Request{Token: token})

	if reply.Status != http.StatusInternalServerError || reply.Body != (wire.FailureView{Error: "internal error"}) {
		t.Fatalf("a driver Message must not Leave, reply=%v", reply)
	}
}
