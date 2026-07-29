package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

// PeopleAPI Holds the Cast every Story Needs.
type PeopleAPI struct {
	Store PeopleStore
}

// ListPeopleRecords Tells the whole Roster.
func (a PeopleAPI) ListPeopleRecords(w http.ResponseWriter, r *http.Request) {
	people, err := a.Store.SelectPeopleRows()
	if err != nil {
		WriteFailureReply(w, err)
		return
	}

	WriteSuccessReply(w, http.StatusOK, RenderPeopleViews(people))
}

// ShowPersonRecord Tells the Story of one Person.
func (a PeopleAPI) ShowPersonRecord(w http.ResponseWriter, r *http.Request) {
	id, err := ReadPersonNumber(r)
	if err != nil {
		WriteFailureReply(w, err)
		return
	}

	person, err := a.Store.SelectPersonRow(id)

	WriteRecordReply(w, http.StatusOK, person, err)
}

// AddPersonRecord Starts a new Story.
func (a PeopleAPI) AddPersonRecord(w http.ResponseWriter, r *http.Request) {
	person, err := ReadPersonBody(r, 0)
	if err != nil {
		WriteFailureReply(w, err)
		return
	}

	stored, err := a.Store.InsertPersonRow(person)

	WriteRecordReply(w, http.StatusCreated, stored, err)
}

// SavePersonRecord Rewrites a Story that already Exists.
func (a PeopleAPI) SavePersonRecord(w http.ResponseWriter, r *http.Request) {
	id, err := ReadPersonNumber(r)
	if err != nil {
		WriteFailureReply(w, err)
		return
	}

	person, err := ReadPersonBody(r, id)
	if err != nil {
		WriteFailureReply(w, err)
		return
	}

	stored, err := a.Store.UpdatePersonRow(person)

	WriteRecordReply(w, http.StatusOK, stored, err)
}

// DropPersonRecord Ends a Story.
func (a PeopleAPI) DropPersonRecord(w http.ResponseWriter, r *http.Request) {
	id, err := ReadPersonNumber(r)
	if err != nil {
		WriteFailureReply(w, err)
		return
	}

	if err := a.Store.DeletePersonRow(id); err != nil {
		WriteFailureReply(w, err)
		return
	}

	WriteSuccessReply(w, http.StatusNoContent, nil)
}

// ReadPersonNumber Pulls the Identity out of the Path.
func ReadPersonNumber(r *http.Request) (PersonID, error) {
	raw, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		return 0, ErrPathIsBroken
	}

	return PersonID(raw), nil
}

// ReadPersonBody Decodes the Wire, then Validates the Business.
func ReadPersonBody(r *http.Request, id PersonID) (Person, error) {
	var body PersonBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return Person{}, ErrBodyIsBroken
	}

	person := body.BuildPersonRecord(id)

	return person, person.CheckPersonRecord()
}

// WriteRecordReply Answers with a Person, or with why there is none.
func WriteRecordReply(w http.ResponseWriter, code int, p Person, err error) {
	if err != nil {
		WriteFailureReply(w, err)
		return
	}

	WriteSuccessReply(w, code, RenderPersonView(p))
}

// WriteSuccessReply Sends the Answer as JSON.
func WriteSuccessReply(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if body != nil {
		_ = json.NewEncoder(w).Encode(body)
	}
}

// WriteFailureReply Maps a Business Failure onto an HTTP Code.
func WriteFailureReply(w http.ResponseWriter, err error) {
	code := http.StatusBadRequest
	if errors.Is(err, ErrPersonUnknown) {
		code = http.StatusNotFound
	}

	log.Printf("❌ %v", err)

	http.Error(w, err.Error(), code)
}
