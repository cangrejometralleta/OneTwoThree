package app

import (
	"log/slog"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/school"
)

// SchoolService Joins the Providers behind each School Use Case.
type SchoolService struct {
	Students StudentStore
	Courses  CourseStore
	Registry RutRegistry
	Notices  EnrollmentNotifier
	Logger   *slog.Logger
}

// ListStudents Reads one Page of the Roll.
func (s SchoolService) ListStudents(page school.Page) ([]school.Student, error) {
	return s.Students.SelectStudentPage(page)
}

// ReadStudent Finds one Enrolment.
func (s SchoolService) ReadStudent(id school.StudentID) (school.Student, error) {
	return s.Students.SelectStudentRow(id)
}

// EnrollStudent Validates and Saves one Enrolment, then Announces it.
// An Announcement that Fails never Fails the Enrolment: the Student Exists.
func (s SchoolService) EnrollStudent(student school.Student) (school.Student, error) {
	if err := s.checkStudent(student); err != nil {
		return school.Student{}, err
	}

	kept, err := s.Students.InsertStudentRow(student)
	if err != nil {
		return school.Student{}, err
	}

	s.announceEnrollment(kept)
	return kept, nil
}

// SaveStudent Validates and Rewrites one Enrolment.
func (s SchoolService) SaveStudent(student school.Student) (school.Student, error) {
	if err := s.checkStudent(student); err != nil {
		return school.Student{}, err
	}

	return s.Students.UpdateStudentRow(student)
}

// DeleteStudent Removes one Enrolment.
func (s SchoolService) DeleteStudent(id school.StudentID) error {
	return s.Students.DeleteStudentRow(id)
}

// ListCourses Reads one Page of the Catalogue.
func (s SchoolService) ListCourses(page school.Page) ([]school.Course, error) {
	return s.Courses.SelectCoursePage(page)
}

// ReadCourse Finds one Course.
func (s SchoolService) ReadCourse(id school.CourseID) (school.Course, error) {
	return s.Courses.SelectCourseRow(id)
}

// CreateCourse Stores one Course.
func (s SchoolService) CreateCourse(course school.Course) (school.Course, error) {
	if err := course.CheckCourseRecord(); err != nil {
		return school.Course{}, err
	}

	return s.Courses.InsertCourseRow(course)
}

func (s SchoolService) checkStudent(student school.Student) error {
	if err := student.CheckStudentRecord(); err != nil {
		return err
	}

	if _, err := s.Courses.SelectCourseRow(student.Course); err != nil {
		return err
	}

	return s.confirmRut(student.Rut)
}

// confirmRut Asks the Registry. An Unknown is never a Yes.
func (s SchoolService) confirmRut(rut school.RUT) error {
	if s.Registry == nil {
		return school.ErrRegistryUnavailable
	}

	registered, err := s.Registry.ConfirmRut(rut)
	if err != nil {
		return school.ErrRegistryUnavailable
	}

	if !registered {
		return school.ErrRutUnregistered
	}

	return nil
}

// announceEnrollment Keeps a Notice Failure inside the Answer: the Caller
// cannot Act on a Notice that did not Leave, so it is Logged and no more.
func (s SchoolService) announceEnrollment(student school.Student) {
	if s.Notices == nil {
		return
	}

	if err := s.Notices.AnnounceEnrollment(student); err != nil {
		s.log().Warn("Enrollment Announcement Failed", "student", student.ID, "error", err)
	}
}

func (s SchoolService) log() *slog.Logger {
	if s.Logger == nil {
		return slog.New(slog.DiscardHandler)
	}

	return s.Logger
}
