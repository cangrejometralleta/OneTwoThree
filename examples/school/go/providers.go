package main

// A Provider is an Interface the Core Declares
// and something outside Fulfils.
// The Core Depends on the Shape, never on the Library behind it.
//
// The Word Collides, and the Collision is worth Knowing.
// Angular and NestJS Call a registered Dependency a Provider.
// JSR-330 Calls a Factory a Provider. Terraform Calls a Plugin one.
// Here it Means the Door an outside System Enters through,
// which the Literature Calls a Port.
// Reference: https://alistair.cockburn.us/hexagonal-architecture/

// StudentStore Keeps Students wherever Students Live.
// Fulfilled by GormSchool.
// Reference: https://gorm.io/docs/
type StudentStore interface {
	InsertStudentRow(s Student) (Student, error)
	SelectStudentRow(id StudentID) (Student, error)
	UpdateStudentRow(s Student) (Student, error)
	DeleteStudentRow(id StudentID) error
	SelectStudentPage(page Page) ([]Student, error)
}

// CourseStore Keeps Courses, and Answers whether one Exists.
// Fulfilled by GormSchool.
// Reference: https://gorm.io/docs/
type CourseStore interface {
	InsertCourseRow(c Course) (Course, error)
	SelectCourseRow(id CourseID) (Course, error)
	SelectCoursePage(page Page) ([]Course, error)
}

// TokenIssuer Mints and Reads a Bearer Token.
// Fulfilled by HmacTokens, so no third Party Enters for this.
// Reference: https://pkg.go.dev/crypto/hmac
type TokenIssuer interface {
	IssueAccessToken(subject string) (string, error)
	ReadAccessToken(token string) (string, error)
}

// Page Carries Pagination without Naming a Database.
type Page struct {
	Number int
	Size   int
}

// CheckPageBounds Refuses a Page the Store cannot Serve.
func (p Page) CheckPageBounds() error {
	if p.Number < 0 || p.Size < 0 {
		return ErrPageIsInvalid
	}

	return nil
}
