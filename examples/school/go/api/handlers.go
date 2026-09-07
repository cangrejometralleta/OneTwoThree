package api

import (
	"net/http"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/app"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/school"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/wire"
)

// SchoolAPI Holds the Cast every Story Needs.
// Three Providers, no Libraries: a Test can Hand it three Fakes.
type SchoolAPI struct {
	Students school.StudentStore
	Courses  school.CourseStore
	Tokens   school.TokenIssuer
}

// DeclareSchoolRoutes is the Libretto.
// Read it once and you Know the whole Service.
func (a SchoolAPI) DeclareSchoolRoutes() []transport.Route {
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
func (a SchoolAPI) Guarded(status int, tell app.Telling) transport.Handler {
	return app.AnswerWith(status, app.RequireProvenCaller(a.Tokens, tell))
}

// MintAccessToken Hands out a Token.
func (a SchoolAPI) MintAccessToken(transport.Request) (any, error) {
	token, err := a.Tokens.IssueAccessToken("student-registry")

	return map[string]string{"token": token}, err
}

// ListStudentRecords Tells one Page of the Roll.
func (a SchoolAPI) ListStudentRecords(req transport.Request) (any, error) {
	page, err := app.ReadPageRequest(req)
	if err != nil {
		return nil, err
	}

	students, err := a.Students.SelectStudentPage(page)
	if err != nil {
		return nil, err
	}

	return school.RenderStudentViews(students), nil
}

// ShowStudentRecord Tells the Story of one Student.
func (a SchoolAPI) ShowStudentRecord(req transport.Request) (any, error) {
	id, err := app.ReadPathNumber(req)
	if err != nil {
		return nil, err
	}

	student, err := a.Students.SelectStudentRow(school.StudentID(id))
	if err != nil {
		return nil, err
	}

	return school.RenderStudentView(student), nil
}

// AddStudentRecord Enrols someone new.
func (a SchoolAPI) AddStudentRecord(req transport.Request) (any, error) {
	student, err := a.ReadStudentBody(req, 0)
	if err != nil {
		return nil, err
	}

	stored, err := a.Students.InsertStudentRow(student)
	if err != nil {
		return nil, err
	}

	return school.RenderStudentView(stored), nil
}

// SaveStudentRecord Rewrites an Enrolment that already Exists.
func (a SchoolAPI) SaveStudentRecord(req transport.Request) (any, error) {
	id, err := app.ReadPathNumber(req)
	if err != nil {
		return nil, err
	}

	student, err := a.ReadStudentBody(req, school.StudentID(id))
	if err != nil {
		return nil, err
	}

	stored, err := a.Students.UpdateStudentRow(student)
	if err != nil {
		return nil, err
	}

	return school.RenderStudentView(stored), nil
}

// DropStudentRecord Ends an Enrolment.
func (a SchoolAPI) DropStudentRecord(req transport.Request) (any, error) {
	id, err := app.ReadPathNumber(req)
	if err != nil {
		return nil, err
	}

	return nil, a.Students.DeleteStudentRow(school.StudentID(id))
}

// ListCourseRecords Tells one Page of the Catalogue.
func (a SchoolAPI) ListCourseRecords(req transport.Request) (any, error) {
	page, err := app.ReadPageRequest(req)
	if err != nil {
		return nil, err
	}

	courses, err := a.Courses.SelectCoursePage(page)
	if err != nil {
		return nil, err
	}

	return school.RenderCourseViews(courses), nil
}

// ShowCourseRecord Tells the Story of one Course.
func (a SchoolAPI) ShowCourseRecord(req transport.Request) (any, error) {
	id, err := app.ReadPathNumber(req)
	if err != nil {
		return nil, err
	}

	course, err := a.Courses.SelectCourseRow(school.CourseID(id))
	if err != nil {
		return nil, err
	}

	return school.RenderCourseView(course), nil
}

// AddCourseRecord Opens a new Course.
func (a SchoolAPI) AddCourseRecord(req transport.Request) (any, error) {
	body, err := app.ReadJSONBody[wire.CourseBody](req)
	if err != nil {
		return nil, err
	}

	stored, err := a.Courses.InsertCourseRow(school.BuildCourseRecord(body, 0))
	if err != nil {
		return nil, err
	}

	return school.RenderCourseView(stored), nil
}

// ReadStudentBody Decodes the Wire, Validates it, then Checks the Course.
func (a SchoolAPI) ReadStudentBody(req transport.Request, id school.StudentID) (school.Student, error) {
	body, err := app.ReadJSONBody[wire.StudentBody](req)
	if err != nil {
		return school.Student{}, err
	}

	student := school.BuildStudentRecord(body, id)
	if err := student.CheckStudentRecord(); err != nil {
		return school.Student{}, err
	}

	_, err = a.Courses.SelectCourseRow(student.Course)

	return student, err
}
