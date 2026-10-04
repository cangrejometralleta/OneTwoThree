package app

import (
	"errors"
	"testing"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/school"
)

// shelf Stands in for the Store: one Course, and every Student it Keeps.
type shelf struct {
	students []school.Student
}

func (s *shelf) InsertStudentRow(student school.Student) (school.Student, error) {
	student.ID = school.StudentID(len(s.students) + 1)
	s.students = append(s.students, student)
	return student, nil
}

func (s *shelf) SelectStudentRow(school.StudentID) (school.Student, error) {
	return school.Student{}, school.ErrStudentUnknown
}

func (s *shelf) UpdateStudentRow(student school.Student) (school.Student, error) {
	return student, nil
}

func (s *shelf) DeleteStudentRow(school.StudentID) error { return nil }

func (s *shelf) SelectStudentPage(school.Page) ([]school.Student, error) { return nil, nil }

func (s *shelf) SelectCourseRow(id school.CourseID) (school.Course, error) {
	if id != 1 {
		return school.Course{}, school.ErrCourseUnknown
	}
	return school.Course{ID: 1, Code: "MAT-101", Name: "Algebra"}, nil
}

func (s *shelf) InsertCourseRow(course school.Course) (school.Course, error) { return course, nil }

func (s *shelf) SelectCoursePage(school.Page) ([]school.Course, error) { return nil, nil }

// office Answers both Questions, the Way one Adapter may Fill two Ports.
type office struct {
	registered bool
	down       error
	noticeErr  error
	asked      []school.RUT
	announced  []school.Student
}

func (o *office) ConfirmRut(rut school.RUT) (bool, error) {
	o.asked = append(o.asked, rut)
	return o.registered, o.down
}

func (o *office) AnnounceEnrollment(student school.Student) error {
	o.announced = append(o.announced, student)
	return o.noticeErr
}

func buildAdmission(front *office) (SchoolService, *shelf) {
	books := &shelf{}
	return SchoolService{Students: books, Courses: books, Registry: front, Notices: front}, books
}

func readyStudent() school.Student {
	return school.Student{Rut: "11111111-1", Name: "Ana", Age: 20, Course: 1}
}

func TestEnrollStudentAnnouncesOnceTheStoreKeeps(t *testing.T) {
	front := &office{registered: true}
	service, books := buildAdmission(front)

	kept, err := service.EnrollStudent(readyStudent())

	if err != nil || len(books.students) != 1 || len(front.announced) != 1 || front.announced[0].ID != kept.ID {
		t.Fatalf("kept=%v err=%v stored=%d announced=%d", kept, err, len(books.students), len(front.announced))
	}
}

func TestEnrollStudentRefusesWhatTheRegistryDenies(t *testing.T) {
	front := &office{registered: false}
	service, books := buildAdmission(front)

	_, err := service.EnrollStudent(readyStudent())

	if !errors.Is(err, school.ErrRutUnregistered) || len(books.students) != 0 || len(front.announced) != 0 {
		t.Fatalf("err=%v stored=%d announced=%d", err, len(books.students), len(front.announced))
	}
}

func TestEnrollStudentNeverTreatsAnUnknownAsAYes(t *testing.T) {
	front := &office{registered: true, down: errors.New("socket closed")}
	service, books := buildAdmission(front)

	_, err := service.EnrollStudent(readyStudent())

	if !errors.Is(err, school.ErrRegistryUnavailable) || len(books.students) != 0 {
		t.Fatalf("err=%v stored=%d", err, len(books.students))
	}

	service.Registry = nil
	if _, err := service.EnrollStudent(readyStudent()); !errors.Is(err, school.ErrRegistryUnavailable) {
		t.Fatalf("a missing Registry must not Pass, err=%v", err)
	}
}

func TestEnrollStudentKeepsAFailedNoticeInsideTheAnswer(t *testing.T) {
	front := &office{registered: true, noticeErr: errors.New("mail is down")}
	service, books := buildAdmission(front)

	_, err := service.EnrollStudent(readyStudent())

	if err != nil || len(books.students) != 1 {
		t.Fatalf("a Notice that did not Leave must not Fail the Enrolment, err=%v stored=%d", err, len(books.students))
	}
}

func TestTheRegistryIsAskedAfterTheFormAndTheCourse(t *testing.T) {
	front := &office{registered: true}
	service, _ := buildAdmission(front)

	invalid := readyStudent()
	invalid.Name = ""
	unknownCourse := readyStudent()
	unknownCourse.Course = 9

	_, formErr := service.EnrollStudent(invalid)
	_, courseErr := service.EnrollStudent(unknownCourse)

	if !errors.Is(formErr, school.ErrNameIsEmpty) || !errors.Is(courseErr, school.ErrCourseUnknown) || len(front.asked) != 0 {
		t.Fatalf("form=%v course=%v asked=%d", formErr, courseErr, len(front.asked))
	}
}

func TestSaveStudentAsksTheRegistryAndAnnouncesNothing(t *testing.T) {
	front := &office{registered: false}
	service, _ := buildAdmission(front)

	if _, err := service.SaveStudent(readyStudent()); !errors.Is(err, school.ErrRutUnregistered) {
		t.Fatalf("a rewrite must ask the Registry too, err=%v", err)
	}

	front.registered = true
	if _, err := service.SaveStudent(readyStudent()); err != nil || len(front.announced) != 0 {
		t.Fatalf("a rewrite announces nothing, err=%v announced=%d", err, len(front.announced))
	}
}
