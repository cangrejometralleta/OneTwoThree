package app

import "github.com/cangrejometralleta/OneTwoThree/examples/school/go/school"

// StudentStore Keeps Students wherever Students Live.
type StudentStore interface {
	InsertStudentRow(student school.Student) (school.Student, error)
	SelectStudentRow(id school.StudentID) (school.Student, error)
	UpdateStudentRow(student school.Student) (school.Student, error)
	DeleteStudentRow(id school.StudentID) error
	SelectStudentPage(page school.Page) ([]school.Student, error)
}

// CourseStore Keeps Courses, and Answers whether one Exists.
type CourseStore interface {
	InsertCourseRow(course school.Course) (school.Course, error)
	SelectCourseRow(id school.CourseID) (school.Course, error)
	SelectCoursePage(page school.Page) ([]school.Course, error)
}

// TokenIssuer Mints a Bearer Token and Reads one for the Guard.
type TokenIssuer interface {
	IssueAccessToken(subject string) (string, error)
	CallerReader
}
