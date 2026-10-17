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
	next     school.StudentID
}

func BuildFakeSchool() *FakeSchool {
	return &FakeSchool{
		students: map[school.StudentID]school.Student{},
		courses:  map[school.CourseID]school.Course{1: {ID: 1, Code: "MAT-101", Name: "Algebra"}},
		next:     1,
	}
}

func (f *FakeSchool) InsertStudentRow(s school.Student) (school.Student, error) {
	for _, enrolled := range f.students {
		if enrolled.Rut == s.Rut {
			return school.Student{}, school.ErrRutTaken
		}
	}
	s.ID = f.next
	f.next++
	f.students[s.ID] = s

	return s, nil
}

func (f *FakeSchool) SelectStudentRow(id school.StudentID) (school.Student, bool, error) {
	student, found := f.students[id]

	return student, found, nil
}

func (f *FakeSchool) UpdateStudentRow(s school.Student) (bool, error) {
	if _, found := f.students[s.ID]; !found {
		return false, nil
	}
	f.students[s.ID] = s

	return true, nil
}

func (f *FakeSchool) DeleteStudentRow(id school.StudentID) (bool, error) {
	if _, found := f.students[id]; !found {
		return false, nil
	}
	delete(f.students, id)

	return true, nil
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

func (f *FakeSchool) SelectCourseRow(id school.CourseID) (school.Course, bool, error) {
	course, found := f.courses[id]

	return course, found, nil
}

func (f *FakeSchool) SelectCoursePage(school.Page) ([]school.Course, error) {
	catalogue := make([]school.Course, 0, len(f.courses))
	for _, course := range f.courses {
		catalogue = append(catalogue, course)
	}
	sort.Slice(catalogue, func(a, b int) bool { return catalogue[a].ID < catalogue[b].ID })

	return catalogue, nil
}

// FakeOffice Stands in for the Registry and the Notice, and Remembers what it was Asked.
type FakeOffice struct {
	Denies bool
	Silent bool
	Muted  bool
	Asked  []school.RUT
}

func (o *FakeOffice) ConfirmRut(rut school.RUT) (bool, error) {
	o.Asked = append(o.Asked, rut)
	if o.Silent {
		return false, errors.New("registry socket closed")
	}

	return !o.Denies, nil
}

func (o *FakeOffice) AnnounceEnrollment(school.Student) error {
	if o.Muted {
		return errors.New("mail is down")
	}

	return nil
}

// BuildTestingSchool Hands the Door its Fakes.
func BuildTestingSchool() (Handler, *FakeOffice) {
	fake := BuildFakeSchool()
	front := &FakeOffice{}

	return Handler{
		School: app.SchoolService{Rules: school.SchoolRules{MinimumAgeYears: 18}, Students: fake, Courses: fake,
			Registry: front, Notices: front},
		Tokens: tokens.BuildAccessTokens("test", 60),
	}, front
}

// CallRoute Arrives the Way a Client Arrives: through the Libretto, the Answer and the Guard.
func CallRoute(t *testing.T, api Handler, method, pattern string, req transport.Request) transport.Response {
	t.Helper()

	token, err := api.Tokens.IssueAccessToken(CallerSubject)
	if err != nil {
		t.Fatal(err)
	}
	req.Headers = map[string]string{"authorization": "Bearer " + token}

	for _, route := range api.DeclareSchoolRoutes().Routes {
		if route.Method == method && route.Pattern == pattern {
			return route.Handle(req)
		}
	}

	t.Fatalf("no Route Answers %s %s", method, pattern)

	return transport.Response{}
}

const enrolBody = `{"rut":"12345678-5","name":"Ada","age":20,"courseId":1}`

func TestEachRouteAnswersWithItsDeclaredStatus(t *testing.T) {
	api, _ := BuildTestingSchool()
	one := transport.Request{Path: map[string]string{"id": "1"}}
	rewrite := transport.Request{
		Path: map[string]string{"id": "1"},
		Body: []byte(`{"rut":"12345678-5","name":"Ada Lovelace","age":21,"courseId":1}`),
	}

	cases := []struct {
		story   string
		method  string
		pattern string
		req     transport.Request
		status  int
	}{
		{"a Token is Minted", "POST", "/token", transport.Request{}, http.StatusCreated},
		{"someone Enrols", "POST", "/students", transport.Request{Body: []byte(enrolBody)}, http.StatusCreated},
		{"the Roll is Read", "GET", "/students", transport.Request{}, http.StatusOK},
		{"one Student is Read", "GET", "/students/{id}", one, http.StatusOK},
		{"an Enrolment is Rewritten", "PUT", "/students/{id}", rewrite, http.StatusOK},
		{"the Catalogue is Read", "GET", "/courses", transport.Request{}, http.StatusOK},
		{"a Course Opens", "POST", "/courses", transport.Request{Body: []byte(`{"code":"FIS-201","name":"Physics"}`)}, http.StatusCreated},
		{"one Course is Read", "GET", "/courses/{id}", one, http.StatusOK},
		{"an Enrolment Ends", "DELETE", "/students/{id}", one, http.StatusNoContent},
	}

	for _, test := range cases {
		if reply := CallRoute(t, api, test.method, test.pattern, test.req); reply.Status != test.status {
			t.Errorf("%s: wanted %d, got %d: %v", test.story, test.status, reply.Status, reply.Body)
		}
	}
}

func TestTheLibrettoDeclaresExactlyTheNineRoutesOfTheSpec(t *testing.T) {
	api, _ := BuildTestingSchool()

	if routes := api.DeclareSchoolRoutes().Routes; len(routes) != 9 {
		t.Fatalf("the Spec Names nine Routes, the Libretto Declares %d", len(routes))
	}
}

func TestARewriteChangesTheStudentAndADeleteLeavesNoBody(t *testing.T) {
	api, _ := BuildTestingSchool()
	CallRoute(t, api, "POST", "/students", transport.Request{Body: []byte(enrolBody)})

	rewritten := CallRoute(t, api, "PUT", "/students/{id}", transport.Request{
		Path: map[string]string{"id": "1"},
		Body: []byte(`{"rut":"12345678-5","name":"Ada Lovelace","age":21,"courseId":1}`),
	})
	if rewritten.Body != (wire.StudentView{ID: 1, Rut: "12345678-5", Name: "Ada Lovelace", Age: 21, Course: 1}) {
		t.Fatalf("the Rewrite must Answer the Student as it Stands, got %v", rewritten.Body)
	}

	dropped := CallRoute(t, api, "DELETE", "/students/{id}", transport.Request{Path: map[string]string{"id": "1"}})
	if dropped.Status != http.StatusNoContent || dropped.Body != nil {
		t.Fatalf("a Delete Answers 204 and no Body, got %d %v", dropped.Status, dropped.Body)
	}
}

func TestTheViewEchoesTheRutInLowerCase(t *testing.T) {
	api, front := BuildTestingSchool()
	CallRoute(t, api, "POST", "/students", transport.Request{Body: []byte(`{"rut":"10000013-K","name":"Ada","age":20,"courseId":1}`)})

	shown := CallRoute(t, api, "GET", "/students/{id}", transport.Request{Path: map[string]string{"id": "1"}})
	if shown.Body != (wire.StudentView{ID: 1, Rut: "10000013-k", Name: "Ada", Age: 20, Course: 1}) || front.Asked[0] != "10000013-k" {
		t.Fatalf("shown=%v asked=%v", shown.Body, front.Asked)
	}
}

func TestAnEmptyListAnswersAnEmptyArray(t *testing.T) {
	api, _ := BuildTestingSchool()

	views, ok := CallRoute(t, api, "GET", "/students", transport.Request{}).Body.([]wire.StudentView)
	if !ok || views == nil || len(views) != 0 {
		t.Fatalf("an empty Roll is [], never null, got %#v", views)
	}
}

func TestEnrolmentSurvivesAFailedAnnouncement(t *testing.T) {
	api, front := BuildTestingSchool()
	front.Muted = true

	if reply := CallRoute(t, api, "POST", "/students", transport.Request{Body: []byte(enrolBody)}); reply.Status != http.StatusCreated {
		t.Fatalf("wanted %d, got %d: %v", http.StatusCreated, reply.Status, reply.Body)
	}
}

func TestEveryRouteButTheTokenRefusesAnUnnamedCaller(t *testing.T) {
	api, _ := BuildTestingSchool()

	for _, route := range api.DeclareSchoolRoutes().Routes {
		if route.Pattern == "/token" {
			continue
		}
		reply := route.Handle(transport.Request{Path: map[string]string{"id": "1"}})
		if reply.Status != http.StatusUnauthorized || reply.Body != (wire.FailureView{Error: "token is Invalid or Expired"}) {
			t.Errorf("%s %s must Refuse an unnamed caller, got %d %v", route.Method, route.Pattern, reply.Status, reply.Body)
		}
	}
}

func TestAnUnknownRouteAnswersBeforeTheGuard(t *testing.T) {
	api, _ := BuildTestingSchool()

	if reply := api.DeclareSchoolRoutes().Unknown(transport.Request{}); reply != (transport.Response{Status: http.StatusNotFound, Body: wire.FailureView{Error: "route not Found"}}) {
		t.Fatalf("reply=%v", reply)
	}
}

func TestAGuardedHandlerReceivesTheNamedCaller(t *testing.T) {
	api, _ := BuildTestingSchool()
	var seen string
	guarded := api.Guarded(http.StatusOK, func(req transport.Request) (any, error) {
		seen = req.Caller

		return "told", nil
	})
	token, _ := api.Tokens.IssueAccessToken(CallerSubject)

	if reply := guarded(transport.Request{Headers: map[string]string{"authorization": "Bearer " + token}}); reply.Status != http.StatusOK || seen != CallerSubject {
		t.Fatalf("reply=%v seen=%q", reply, seen)
	}
}

func TestAnUncontrolledFailureAnswersFiveHundredAndSaysNothingOfTheDriver(t *testing.T) {
	api, _ := BuildTestingSchool()
	broken := api.Guarded(http.StatusOK, func(transport.Request) (any, error) {
		return nil, errors.New("the driver Broke at /var/db")
	})
	token, _ := api.Tokens.IssueAccessToken(CallerSubject)

	reply := broken(transport.Request{Headers: map[string]string{"authorization": "Bearer " + token}})

	if reply.Status != http.StatusInternalServerError || reply.Body != (wire.FailureView{Error: "internal error"}) {
		t.Fatalf("a driver Message must not Leave, reply=%v", reply)
	}
}
