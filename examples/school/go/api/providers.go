package api

import (
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/app"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/school"
)

// A Provider is an Interface the Consumer Declares
// and something outside Fulfils.
// The Handlers Call these, so the Handlers Declare them, small and here.
// main.go is the one Place that Joins them to what Fulfils them.
//
// The Word Collides, and the Collision is worth Knowing.
// Angular and NestJS Call a registered Dependency a Provider.
// JSR-330 Calls a Factory a Provider. Terraform Calls a Plugin one.
// Here it Means the Door an outside System Enters through,
// which the Literature Calls a Port.
// Reference: https://alistair.cockburn.us/hexagonal-architecture/

// StudentStore Keeps Students wherever Students Live.
// Fulfilled by store.School.
// Reference: https://gorm.io/docs/
type StudentStore interface {
	InsertStudentRow(s school.Student) (school.Student, error)
	SelectStudentRow(id school.StudentID) (school.Student, error)
	UpdateStudentRow(s school.Student) (school.Student, error)
	DeleteStudentRow(id school.StudentID) error
	SelectStudentPage(page school.Page) ([]school.Student, error)
}

// CourseStore Keeps Courses, and Answers whether one Exists.
// Fulfilled by store.School.
// Reference: https://gorm.io/docs/
type CourseStore interface {
	InsertCourseRow(c school.Course) (school.Course, error)
	SelectCourseRow(id school.CourseID) (school.Course, error)
	SelectCoursePage(page school.Page) ([]school.Course, error)
}

// TokenIssuer Mints a Bearer Token, and Reads one for the Guard.
// The Reading Half is Declared by app, where the Guard Calls it.
// Fulfilled by tokens.AccessTokens, so no third Party Enters for this.
// Reference: https://pkg.go.dev/crypto/hmac
type TokenIssuer interface {
	IssueAccessToken(subject string) (string, error)
	app.CallerReader
}
