package api

import (
	"errors"
	"net/http"
	"testing"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/school"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/tokens"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
)

// FakeSchool Stands in for GORM.
// No Database, no Framework, no Port: the Providers Allow it.
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

// InsertStudentRow Keeps the RUT unique, the Way a real Index would.
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
	students := make([]school.Student, 0, len(f.students))
	for _, student := range f.students {
		students = append(students, student)
	}
	return students, nil
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
	courses := make([]school.Course, 0, len(f.courses))
	for _, course := range f.courses {
		courses = append(courses, course)
	}
	return courses, nil
}

// BuildTestingSchool Hands the API three Fakes.
func BuildTestingSchool() SchoolAPI {
	fake := BuildFakeSchool()
	issuer := tokens.BuildAccessTokens("test", 60)

	return SchoolAPI{Students: fake, Courses: fake, Tokens: issuer}
}

// CallRoute Arrives the Way a Client Arrives: through the Libretto.
// It Crosses AnswerWith and the Guard, so a Test Sees the real Status.
func CallRoute(t *testing.T, api SchoolAPI, method, pattern string, req transport.Request) transport.Response {
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
		{"one Student is Read", "GET", "/students/{id}", transport.Request{Path: map[string]string{"id": "1"}}, http.StatusOK},
		{"an Enrolment Ends", "DELETE", "/students/{id}", transport.Request{Path: map[string]string{"id": "1"}}, http.StatusNoContent},
	}

	for _, test := range cases {
		t.Run(test.story, func(t *testing.T) {
			if reply := CallRoute(t, api, test.method, test.pattern, test.req); reply.Status != test.status {
				t.Fatalf("wanted %d, got %d: %v", test.status, reply.Status, reply.Body)
			}
		})
	}
}

// Every Route but the Token Names its Caller first.
func TestEveryRouteButTheTokenRefusesAnUnnamedCaller(t *testing.T) {
	api := BuildTestingSchool()

	for _, route := range api.DeclareSchoolRoutes() {
		reply := route.Handle(transport.Request{Path: map[string]string{"id": "1"}})

		wanted := http.StatusUnauthorized
		if route.Pattern == "/token" {
			wanted = http.StatusCreated
		}

		if reply.Status != wanted {
			t.Errorf("%s %s: wanted %d, got %d", route.Method, route.Pattern, wanted, reply.Status)
		}
	}
}

// A Guarded Handler Reads a Caller it never had to Check.
func TestGuardedHandlerReceivesTheNamedCaller(t *testing.T) {
	api := BuildTestingSchool()
	seen := ""

	guarded := api.Guarded(http.StatusOK, func(req transport.Request) (any, error) {
		seen = req.Caller

		return nil, nil
	})

	token, _ := api.Tokens.IssueAccessToken("student-registry")
	guarded(transport.Request{Token: token})

	if seen != "student-registry" {
		t.Fatalf("wanted the Subject, got %q", seen)
	}
}

// An uncontrolled Store Failure never Becomes the Caller's Fault.
func TestUncontrolledFailureAnswersFiveHundred(t *testing.T) {
	api := BuildTestingSchool()
	api.Students = BrokenSchool{FakeSchool: BuildFakeSchool()}

	reply := CallRoute(t, api, "GET", "/students", transport.Request{})

	if reply.Status != http.StatusInternalServerError {
		t.Fatalf("wanted 500, got %d", reply.Status)
	}
}

// BrokenSchool Fails the Way a Driver Fails: with no Status to Offer.
type BrokenSchool struct {
	*FakeSchool
}

func (BrokenSchool) SelectStudentPage(school.Page) ([]school.Student, error) {
	return nil, errors.New("the driver Closed the Connection")
}
