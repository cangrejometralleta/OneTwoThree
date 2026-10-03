package handler

import (
	"net/http"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/app"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/school"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/wire"
)

// Handler Holds the Cast every Story Needs.
// Two Dependencies, no Libraries: a Test can Hand it Fakes.
type Handler struct {
	School app.SchoolService
	Tokens app.TokenIssuer
}

// DeclareSchoolRoutes is the Libretto.
// Read it once and you Know the whole Service.
func (a Handler) DeclareSchoolRoutes() []transport.Route {
	return []transport.Route{
		{Method: "POST", Pattern: "/token", Handle: app.AnswerWith(http.StatusCreated, a.MintAccessToken)},

		{Method: "GET", Pattern: "/students", Handle: a.Guarded(http.StatusOK, a.ListStudentRecords)},
		{Method: "POST", Pattern: "/students", Handle: a.Guarded(http.StatusCreated, a.AddStudentRecord)},
		{Method: "GET", Pattern: "/students/{id}", Handle: a.Guarded(http.StatusOK, a.ShowStudentRecord)},
		{Method: "PUT", Pattern: "/students/{id}", Handle: a.Guarded(http.StatusOK, a.SaveStudentRecord)},
		{Method: "DELETE", Pattern: "/students/{id}", Handle: a.Guarded(http.StatusNoContent, a.DropStudentRecord)},

		{Method: "GET", Pattern: "/courses", Handle: a.Guarded(http.StatusOK, a.ListCourseRecords)},
		{Method: "POST", Pattern: "/courses", Handle: a.Guarded(http.StatusCreated, a.AddCourseRecord)},
		{Method: "GET", Pattern: "/courses/{id}", Handle: a.Guarded(http.StatusOK, a.ShowCourseRecord)},
	}
}

// Guarded Says the same two Things about a Route every time:
// Name the Caller first, then Answer with this Status.
func (a Handler) Guarded(status int, tell app.Telling) transport.Handler {
	return app.AnswerWith(status, app.RequireProvenCaller(a.Tokens, tell))
}

// MintAccessToken Hands out a Token.
func (a Handler) MintAccessToken(transport.Request) (any, error) {
	token, err := a.Tokens.IssueAccessToken("student-registry")

	return map[string]string{"token": token}, err
}

// ListStudentRecords Tells one Page of the Roll.
func (a Handler) ListStudentRecords(req transport.Request) (any, error) {
	page, err := app.ReadPageRequest(req)
	if err != nil {
		return nil, err
	}

	students, err := a.School.ListStudents(page)
	if err != nil {
		return nil, err
	}

	return renderStudentViews(students), nil
}

// ShowStudentRecord Tells the Story of one Student.
func (a Handler) ShowStudentRecord(req transport.Request) (any, error) {
	id, err := app.ReadPathNumber(req)
	if err != nil {
		return nil, err
	}

	student, err := a.School.ReadStudent(school.StudentID(id))
	if err != nil {
		return nil, err
	}

	return renderStudentView(student), nil
}

// AddStudentRecord Enrols someone new.
func (a Handler) AddStudentRecord(req transport.Request) (any, error) {
	student, err := a.readStudentBody(req, 0)
	if err != nil {
		return nil, err
	}

	stored, err := a.School.EnrollStudent(student)
	if err != nil {
		return nil, err
	}

	return renderStudentView(stored), nil
}

// SaveStudentRecord Rewrites an Enrolment that already Exists.
func (a Handler) SaveStudentRecord(req transport.Request) (any, error) {
	id, err := app.ReadPathNumber(req)
	if err != nil {
		return nil, err
	}

	student, err := a.readStudentBody(req, school.StudentID(id))
	if err != nil {
		return nil, err
	}

	stored, err := a.School.SaveStudent(student)
	if err != nil {
		return nil, err
	}

	return renderStudentView(stored), nil
}

// DropStudentRecord Ends an Enrolment.
func (a Handler) DropStudentRecord(req transport.Request) (any, error) {
	id, err := app.ReadPathNumber(req)
	if err != nil {
		return nil, err
	}

	return nil, a.School.DeleteStudent(school.StudentID(id))
}

// ListCourseRecords Tells one Page of the Catalogue.
func (a Handler) ListCourseRecords(req transport.Request) (any, error) {
	page, err := app.ReadPageRequest(req)
	if err != nil {
		return nil, err
	}

	courses, err := a.School.ListCourses(page)
	if err != nil {
		return nil, err
	}

	return renderCourseViews(courses), nil
}

// ShowCourseRecord Tells the Story of one Course.
func (a Handler) ShowCourseRecord(req transport.Request) (any, error) {
	id, err := app.ReadPathNumber(req)
	if err != nil {
		return nil, err
	}

	course, err := a.School.ReadCourse(school.CourseID(id))
	if err != nil {
		return nil, err
	}

	return renderCourseView(course), nil
}

// AddCourseRecord Opens a new Course.
func (a Handler) AddCourseRecord(req transport.Request) (any, error) {
	body, err := app.ReadJSONBody[wire.CourseBody](req)
	if err != nil {
		return nil, err
	}

	stored, err := a.School.CreateCourse(buildCourseRecord(body, 0))
	if err != nil {
		return nil, err
	}

	return renderCourseView(stored), nil
}

// readStudentBody Decodes the Wire into the Business Shape.
func (a Handler) readStudentBody(req transport.Request, id school.StudentID) (school.Student, error) {
	body, err := app.ReadJSONBody[wire.StudentBody](req)
	if err != nil {
		return school.Student{}, err
	}

	return buildStudentRecord(body, id), nil
}
