package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// SchoolAPI Holds the Cast every Story Needs.
// Three Providers, no Libraries: a Test can Hand it three Fakes.
type SchoolAPI struct {
	Students StudentStore
	Courses  CourseStore
	Tokens   TokenIssuer
}

// DeclareSchoolRoutes is the Libretto.
// Read it once and you Know the whole Service.
func (a SchoolAPI) DeclareSchoolRoutes() []Route {
	return []Route{
		{"POST", "/token", a.MintAccessToken},

		{"GET", "/students", a.ListStudentRecords},
		{"POST", "/students", a.AddStudentRecord},
		{"GET", "/students/{id}", a.ShowStudentRecord},
		{"PUT", "/students/{id}", a.SaveStudentRecord},
		{"DELETE", "/students/{id}", a.DropStudentRecord},

		{"GET", "/courses", a.ListCourseRecords},
		{"POST", "/courses", a.AddCourseRecord},
		{"GET", "/courses/{id}", a.ShowCourseRecord},
	}
}

// MintAccessToken Hands out a Token, and Guards nothing else.
func (a SchoolAPI) MintAccessToken(req Request) Response {
	token, err := a.Tokens.IssueAccessToken("student-registry")
	if err != nil {
		return BuildFailureReply(err)
	}

	return Response{http.StatusCreated, map[string]string{"token": token}}
}

// ListStudentRecords Tells one Page of the Roll.
func (a SchoolAPI) ListStudentRecords(req Request) Response {
	page, err := ReadPageRequest(req)
	if err != nil {
		return BuildFailureReply(err)
	}

	students, err := a.Students.SelectStudentPage(page)
	if err != nil {
		return BuildFailureReply(err)
	}

	return Response{http.StatusOK, RenderStudentViews(students)}
}

// ShowStudentRecord Tells the Story of one Student.
func (a SchoolAPI) ShowStudentRecord(req Request) Response {
	id, err := ReadPathNumber(req)
	if err != nil {
		return BuildFailureReply(err)
	}

	student, err := a.Students.SelectStudentRow(StudentID(id))

	return BuildStudentReply(http.StatusOK, student, err)
}

// AddStudentRecord Enrols someone new.
func (a SchoolAPI) AddStudentRecord(req Request) Response {
	student, err := a.ReadStudentBody(req, 0)
	if err != nil {
		return BuildFailureReply(err)
	}

	stored, err := a.Students.InsertStudentRow(student)

	return BuildStudentReply(http.StatusCreated, stored, err)
}

// SaveStudentRecord Rewrites an Enrolment that already Exists.
func (a SchoolAPI) SaveStudentRecord(req Request) Response {
	id, err := ReadPathNumber(req)
	if err != nil {
		return BuildFailureReply(err)
	}

	student, err := a.ReadStudentBody(req, StudentID(id))
	if err != nil {
		return BuildFailureReply(err)
	}

	stored, err := a.Students.UpdateStudentRow(student)

	return BuildStudentReply(http.StatusOK, stored, err)
}

// DropStudentRecord Ends an Enrolment.
func (a SchoolAPI) DropStudentRecord(req Request) Response {
	id, err := ReadPathNumber(req)
	if err != nil {
		return BuildFailureReply(err)
	}

	if err := a.Students.DeleteStudentRow(StudentID(id)); err != nil {
		return BuildFailureReply(err)
	}

	return Response{http.StatusNoContent, nil}
}

// ListCourseRecords Tells one Page of the Catalogue.
func (a SchoolAPI) ListCourseRecords(req Request) Response {
	page, err := ReadPageRequest(req)
	if err != nil {
		return BuildFailureReply(err)
	}

	courses, err := a.Courses.SelectCoursePage(page)
	if err != nil {
		return BuildFailureReply(err)
	}

	return Response{http.StatusOK, RenderCourseViews(courses)}
}

// ShowCourseRecord Tells the Story of one Course.
func (a SchoolAPI) ShowCourseRecord(req Request) Response {
	id, err := ReadPathNumber(req)
	if err != nil {
		return BuildFailureReply(err)
	}

	course, err := a.Courses.SelectCourseRow(CourseID(id))
	if err != nil {
		return BuildFailureReply(err)
	}

	return Response{http.StatusOK, RenderCourseView(course)}
}

// AddCourseRecord Opens a new Course.
func (a SchoolAPI) AddCourseRecord(req Request) Response {
	var body CourseBody
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return BuildFailureReply(ErrBodyIsBroken)
	}

	stored, err := a.Courses.InsertCourseRow(body.BuildCourseRecord(0))
	if err != nil {
		return BuildFailureReply(err)
	}

	return Response{http.StatusCreated, RenderCourseView(stored)}
}

// ReadStudentBody Decodes the Wire, Validates it, then Checks the Course.
func (a SchoolAPI) ReadStudentBody(req Request, id StudentID) (Student, error) {
	var body StudentBody
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return Student{}, ErrBodyIsBroken
	}

	student := body.BuildStudentRecord(id)
	if err := student.CheckStudentRecord(); err != nil {
		return Student{}, err
	}

	_, err := a.Courses.SelectCourseRow(student.Course)

	return student, err
}

// ReadPathNumber Pulls an Identity out of the Path.
func ReadPathNumber(req Request) (uint64, error) {
	raw, err := strconv.ParseUint(req.Path["id"], 10, 64)
	if err != nil {
		return 0, ErrPathIsBroken
	}

	return raw, nil
}

// ReadPageRequest Reads Pagination, Defaulting to the whole Set.
func ReadPageRequest(req Request) (Page, error) {
	number, _ := strconv.Atoi(req.Query["page"])
	size, _ := strconv.Atoi(req.Query["size"])

	page := Page{Number: number, Size: size}

	return page, page.CheckPageBounds()
}

// BuildStudentReply Answers with a Student, or with why there is none.
func BuildStudentReply(code int, s Student, err error) Response {
	if err != nil {
		return BuildFailureReply(err)
	}

	return Response{code, RenderStudentView(s)}
}

// BuildFailureReply Maps a Business Failure onto an HTTP Code.
// The Mapping Lives here once, never inside a Handler.
func BuildFailureReply(err error) Response {
	body := map[string]string{"error": err.Error()}

	if errors.Is(err, ErrStudentUnknown) || errors.Is(err, ErrCourseUnknown) {
		return Response{http.StatusNotFound, body}
	}

	if errors.Is(err, ErrRutTaken) {
		return Response{http.StatusConflict, body}
	}

	return Response{http.StatusBadRequest, body}
}
