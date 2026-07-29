package main

import (
	"net/http"
	"testing"
	"time"
)

// FakeSchool Stands in for GORM.
// No Database, no Framework, no Port: the Providers Allow it.
type FakeSchool struct {
	students map[StudentID]Student
	courses  map[CourseID]Course
}

func BuildFakeSchool() *FakeSchool {
	return &FakeSchool{
		students: map[StudentID]Student{},
		courses:  map[CourseID]Course{1: {ID: 1, Code: "MAT-101", Name: "Algebra"}},
	}
}

func (f *FakeSchool) InsertStudentRow(s Student) (Student, error) {
	s.ID = StudentID(len(f.students) + 1)
	f.students[s.ID] = s
	return s, nil
}

func (f *FakeSchool) SelectStudentRow(id StudentID) (Student, error) {
	student, found := f.students[id]
	if !found {
		return Student{}, ErrStudentUnknown
	}
	return student, nil
}

func (f *FakeSchool) UpdateStudentRow(s Student) (Student, error) {
	if _, found := f.students[s.ID]; !found {
		return Student{}, ErrStudentUnknown
	}
	f.students[s.ID] = s
	return s, nil
}

func (f *FakeSchool) DeleteStudentRow(id StudentID) error {
	if _, found := f.students[id]; !found {
		return ErrStudentUnknown
	}
	delete(f.students, id)
	return nil
}

func (f *FakeSchool) SelectStudentPage(Page) ([]Student, error) {
	students := make([]Student, 0, len(f.students))
	for _, student := range f.students {
		students = append(students, student)
	}
	return students, nil
}

func (f *FakeSchool) InsertCourseRow(c Course) (Course, error) {
	c.ID = CourseID(len(f.courses) + 1)
	f.courses[c.ID] = c
	return c, nil
}

func (f *FakeSchool) SelectCourseRow(id CourseID) (Course, error) {
	course, found := f.courses[id]
	if !found {
		return Course{}, ErrCourseUnknown
	}
	return course, nil
}

func (f *FakeSchool) SelectCoursePage(Page) ([]Course, error) {
	courses := make([]Course, 0, len(f.courses))
	for _, course := range f.courses {
		courses = append(courses, course)
	}
	return courses, nil
}

// BuildTestingSchool Hands the API three Fakes and nothing else.
func BuildTestingSchool() SchoolAPI {
	fake := BuildFakeSchool()
	tokens := HmacTokens{Secret: []byte("test"), Life: time.Minute, Now: time.Now}

	return SchoolAPI{Students: fake, Courses: fake, Tokens: tokens}
}

func TestAddStudentRecordEnrolsSomeone(t *testing.T) {
	reply := BuildTestingSchool().AddStudentRecord(Request{
		Body: []byte(`{"rut":"12345678-5","name":"Ada","age":20,"courseId":1}`),
	})

	if reply.Status != http.StatusCreated {
		t.Fatalf("wanted 201, got %d: %v", reply.Status, reply.Body)
	}
}

func TestAddStudentRecordRefusesBadRut(t *testing.T) {
	reply := BuildTestingSchool().AddStudentRecord(Request{
		Body: []byte(`{"rut":"12345678-9","name":"Ada","age":20,"courseId":1}`),
	})

	if reply.Status != http.StatusBadRequest {
		t.Fatalf("wanted 400, got %d", reply.Status)
	}
}

func TestAddStudentRecordRefusesUnknownCourse(t *testing.T) {
	reply := BuildTestingSchool().AddStudentRecord(Request{
		Body: []byte(`{"rut":"12345678-5","name":"Ada","age":20,"courseId":99}`),
	})

	if reply.Status != http.StatusNotFound {
		t.Fatalf("wanted 404, got %d", reply.Status)
	}
}

func TestShowStudentRecordReportsAbsence(t *testing.T) {
	reply := BuildTestingSchool().ShowStudentRecord(Request{Path: map[string]string{"id": "42"}})

	if reply.Status != http.StatusNotFound {
		t.Fatalf("wanted 404, got %d", reply.Status)
	}
}

func TestListStudentRecordsRefusesNegativePage(t *testing.T) {
	reply := BuildTestingSchool().ListStudentRecords(Request{Query: map[string]string{"page": "-1"}})

	if reply.Status != http.StatusBadRequest {
		t.Fatalf("wanted 400, got %d", reply.Status)
	}
}

// The Comparison that Broke in the Spring Version, Pinned here.
func TestAccessTokenDiesOnlyAfterItsDeadline(t *testing.T) {
	clock := time.Now()
	tokens := HmacTokens{Secret: []byte("s"), Life: time.Minute, Now: func() time.Time { return clock }}

	token, _ := tokens.IssueAccessToken("registry")

	if _, err := tokens.ReadAccessToken(token); err != nil {
		t.Fatalf("a fresh Token Must be accepted, got %v", err)
	}

	clock = clock.Add(2 * time.Minute)
	if _, err := tokens.ReadAccessToken(token); err == nil {
		t.Fatal("an expired Token Must be Refused")
	}
}

func TestAccessTokenRefusesATamperedClaim(t *testing.T) {
	tokens := HmacTokens{Secret: []byte("s"), Life: time.Minute, Now: time.Now}

	token, _ := tokens.IssueAccessToken("registry")

	if _, err := tokens.ReadAccessToken("admin" + token[8:]); err == nil {
		t.Fatal("a rewritten Subject Must be Refused")
	}
}

func TestTranslateRoutePatternFeedsGin(t *testing.T) {
	if got := TranslateRoutePattern("/students/{id}"); got != "/students/:id" {
		t.Fatalf("wanted /students/:id, got %s", got)
	}
}
