package app

import "github.com/cangrejometralleta/OneTwoThree/examples/school/go/school"

// StudentStore Keeps Students wherever Students Live.
// Absence Answers false; only the Store Proves a RUT Taken.
type StudentStore interface {
	InsertStudentRow(student school.Student) (school.Student, error)
	SelectStudentRow(id school.StudentID) (school.Student, bool, error)
	UpdateStudentRow(student school.Student) (bool, error)
	DeleteStudentRow(id school.StudentID) (bool, error)
	SelectStudentPage(page school.Page) ([]school.Student, error)
}

// CourseStore Keeps Courses, and Answers whether one Exists.
type CourseStore interface {
	InsertCourseRow(course school.Course) (school.Course, error)
	SelectCourseRow(id school.CourseID) (school.Course, bool, error)
	SelectCoursePage(page school.Page) ([]school.Course, error)
}

// RutRegistry Confirms whether a RUT Belongs to someone Real.
// A RUT it Denies Answers false; a Registry that cannot Answer returns an error.
type RutRegistry interface {
	ConfirmRut(rut school.RUT) (bool, error)
}

// EnrollmentNotifier Announces a Student the Store has Kept.
type EnrollmentNotifier interface {
	AnnounceEnrollment(student school.Student) error
}

// CallerReader Names the Caller a Token Proves, or Fails.
type CallerReader interface {
	ReadCaller(token string) (string, error)
}

// TokenIssuer Mints a Bearer Token and Reads one for the Guard.
type TokenIssuer interface {
	IssueAccessToken(subject string) (string, error)
	CallerReader
}
