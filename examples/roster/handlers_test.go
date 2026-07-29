package main

import (
	"encoding/json"
	"net/http"
	"testing"
)

// FakePeople Stands in for GORM.
// No Database, no Framework, no Network: the Interface Allows it.
type FakePeople struct {
	rows map[PersonID]Person
}

func (f FakePeople) InsertPersonRow(p Person) (Person, error) {
	p.ID = PersonID(len(f.rows) + 1)
	f.rows[p.ID] = p
	return p, nil
}

func (f FakePeople) SelectPersonRow(id PersonID) (Person, error) {
	person, found := f.rows[id]
	if !found {
		return Person{}, ErrPersonUnknown
	}
	return person, nil
}

func (f FakePeople) UpdatePersonRow(p Person) (Person, error) {
	if _, found := f.rows[p.ID]; !found {
		return Person{}, ErrPersonUnknown
	}
	f.rows[p.ID] = p
	return p, nil
}

func (f FakePeople) DeletePersonRow(id PersonID) error {
	if _, found := f.rows[id]; !found {
		return ErrPersonUnknown
	}
	delete(f.rows, id)
	return nil
}

func (f FakePeople) SelectPeopleRows() ([]Person, error) {
	people := make([]Person, 0, len(f.rows))
	for _, person := range f.rows {
		people = append(people, person)
	}
	return people, nil
}

// BuildTestingRoster Hands the API a Fake and nothing else.
func BuildTestingRoster() RosterAPI {
	return RosterAPI{Store: FakePeople{rows: map[PersonID]Person{}}}
}

func TestAddPersonRecord(t *testing.T) {
	api := BuildTestingRoster()

	reply := api.AddPersonRecord(Request{Body: []byte(`{"name":"Ada","email":"ada@lovelace.dev"}`)})

	if reply.Status != http.StatusCreated {
		t.Fatalf("wanted 201, got %d", reply.Status)
	}
}

func TestAddPersonRecordRefusesEmptyName(t *testing.T) {
	api := BuildTestingRoster()

	reply := api.AddPersonRecord(Request{Body: []byte(`{"name":"","email":"ada@lovelace.dev"}`)})

	if reply.Status != http.StatusBadRequest {
		t.Fatalf("wanted 400, got %d", reply.Status)
	}
}

func TestShowPersonRecordReportsAbsence(t *testing.T) {
	api := BuildTestingRoster()

	reply := api.ShowPersonRecord(Request{Path: map[string]string{"id": "42"}})

	if reply.Status != http.StatusNotFound {
		t.Fatalf("wanted 404, got %d", reply.Status)
	}
}

func TestDeclareRosterRoutesCoversFiveVerbs(t *testing.T) {
	routes := BuildTestingRoster().DeclareRosterRoutes()

	if len(routes) != 5 {
		t.Fatalf("wanted 5 Routes, got %d", len(routes))
	}
}

func TestTranslateRoutePatternFeedsGin(t *testing.T) {
	if got := TranslateRoutePattern("/people/{id}"); got != "/people/:id" {
		t.Fatalf("wanted /people/:id, got %s", got)
	}
}

func TestRenderPersonViewCrossesToTheWire(t *testing.T) {
	view := RenderPersonView(Person{ID: 7, Name: "Ada", Email: "ada@lovelace.dev"})

	raw, _ := json.Marshal(view)
	if string(raw) != `{"id":7,"name":"Ada","email":"ada@lovelace.dev","phone":""}` {
		t.Fatalf("unexpected Wire Shape: %s", raw)
	}
}
