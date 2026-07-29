package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// RosterAPI Holds the Cast every Story Needs.
// It Depends on an Interface, so a Test can Hand it a Fake.
type RosterAPI struct {
	Store PeopleStore
}

// DeclareRosterRoutes is the Libretto.
// Read it once and you Know the whole Service.
func (a RosterAPI) DeclareRosterRoutes() []Route {
	return []Route{
		{"GET", "/people", a.ListPeopleRecords},
		{"POST", "/people", a.AddPersonRecord},
		{"GET", "/people/{id}", a.ShowPersonRecord},
		{"PUT", "/people/{id}", a.SavePersonRecord},
		{"DELETE", "/people/{id}", a.DropPersonRecord},
	}
}

// ListPeopleRecords Tells the whole Roster.
func (a RosterAPI) ListPeopleRecords(req Request) Response {
	people, err := a.Store.SelectPeopleRows()
	if err != nil {
		return BuildFailureReply(err)
	}

	return Response{http.StatusOK, RenderPeopleViews(people)}
}

// ShowPersonRecord Tells the Story of one Person.
func (a RosterAPI) ShowPersonRecord(req Request) Response {
	id, err := ReadPersonNumber(req)
	if err != nil {
		return BuildFailureReply(err)
	}

	person, err := a.Store.SelectPersonRow(id)

	return BuildRecordReply(http.StatusOK, person, err)
}

// AddPersonRecord Starts a new Story.
func (a RosterAPI) AddPersonRecord(req Request) Response {
	person, err := ReadPersonBody(req, 0)
	if err != nil {
		return BuildFailureReply(err)
	}

	stored, err := a.Store.InsertPersonRow(person)

	return BuildRecordReply(http.StatusCreated, stored, err)
}

// SavePersonRecord Rewrites a Story that already Exists.
func (a RosterAPI) SavePersonRecord(req Request) Response {
	id, err := ReadPersonNumber(req)
	if err != nil {
		return BuildFailureReply(err)
	}

	person, err := ReadPersonBody(req, id)
	if err != nil {
		return BuildFailureReply(err)
	}

	stored, err := a.Store.UpdatePersonRow(person)

	return BuildRecordReply(http.StatusOK, stored, err)
}

// DropPersonRecord Ends a Story.
func (a RosterAPI) DropPersonRecord(req Request) Response {
	id, err := ReadPersonNumber(req)
	if err != nil {
		return BuildFailureReply(err)
	}

	if err := a.Store.DeletePersonRow(id); err != nil {
		return BuildFailureReply(err)
	}

	return Response{http.StatusNoContent, nil}
}

// ReadPersonNumber Pulls the Identity out of the Path.
func ReadPersonNumber(req Request) (PersonID, error) {
	raw, err := strconv.ParseUint(req.Path["id"], 10, 64)
	if err != nil {
		return 0, ErrPathIsBroken
	}

	return PersonID(raw), nil
}

// ReadPersonBody Decodes the Wire, then Validates the Business.
func ReadPersonBody(req Request, id PersonID) (Person, error) {
	var body PersonBody
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return Person{}, ErrBodyIsBroken
	}

	person := body.BuildPersonRecord(id)

	return person, person.CheckPersonRecord()
}

// BuildRecordReply Answers with a Person, or with why there is none.
func BuildRecordReply(code int, p Person, err error) Response {
	if err != nil {
		return BuildFailureReply(err)
	}

	return Response{code, RenderPersonView(p)}
}

// BuildFailureReply Maps a Business Failure onto an HTTP Code.
func BuildFailureReply(err error) Response {
	if errors.Is(err, ErrPersonUnknown) {
		return Response{http.StatusNotFound, map[string]string{"error": err.Error()}}
	}

	return Response{http.StatusBadRequest, map[string]string{"error": err.Error()}}
}
