package main

import (
	"errors"

	"gorm.io/gorm"
)

// PersonRow is the Storage Shape.
// The Tags Document the Table, so the Schema Explains itself
// and no separate Migration File Has to.
type PersonRow struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"size:120;not null"`
	Email string `gorm:"size:180;uniqueIndex;not null"`
	Phone string `gorm:"size:40"`
}

// PeopleStore Hides every Query behind a Business Verb.
type PeopleStore interface {
	InsertPersonRow(p Person) (Person, error)
	SelectPersonRow(id PersonID) (Person, error)
	UpdatePersonRow(p Person) (Person, error)
	DeletePersonRow(id PersonID) error
	SelectPeopleRows() ([]Person, error)
}

// GormPeople Speaks Machine, so the Handlers can Speak Business.
type GormPeople struct {
	DB *gorm.DB
}

// InsertPersonRow Writes a new Person and Returns it Numbered.
func (g GormPeople) InsertPersonRow(p Person) (Person, error) {
	row := EncodePersonRow(p)

	if err := g.DB.Create(&row).Error; err != nil {
		return Person{}, TranslateStoreError(err)
	}

	return DecodePersonRow(row), nil
}

// SelectPersonRow Finds one Person or Says why not.
func (g GormPeople) SelectPersonRow(id PersonID) (Person, error) {
	var row PersonRow

	if err := g.DB.First(&row, uint(id)).Error; err != nil {
		return Person{}, TranslateStoreError(err)
	}

	return DecodePersonRow(row), nil
}

// UpdatePersonRow Overwrites an existing Person.
func (g GormPeople) UpdatePersonRow(p Person) (Person, error) {
	row := EncodePersonRow(p)

	result := g.DB.
		Model(&PersonRow{}).
		Where("id = ?", row.ID).
		Updates(row)

	if result.RowsAffected == 0 {
		return Person{}, ErrPersonUnknown
	}

	return p, TranslateStoreError(result.Error)
}

// DeletePersonRow Removes a Person, or Reports an Absence.
func (g GormPeople) DeletePersonRow(id PersonID) error {
	result := g.DB.Delete(&PersonRow{}, uint(id))

	if result.RowsAffected == 0 {
		return ErrPersonUnknown
	}

	return TranslateStoreError(result.Error)
}

// SelectPeopleRows Reads the whole Roster, Ordered by Name.
func (g GormPeople) SelectPeopleRows() ([]Person, error) {
	var rows []PersonRow

	if err := g.DB.Order("name").Find(&rows).Error; err != nil {
		return nil, TranslateStoreError(err)
	}

	return DecodePeopleRows(rows), nil
}

// TranslateStoreError Turns a Driver Failure into a Business one.
func TranslateStoreError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrPersonUnknown
	}

	return err
}

// EncodePersonRow Turns a Business Person into Storage.
func EncodePersonRow(p Person) PersonRow {
	return PersonRow{
		ID:    uint(p.ID),
		Name:  string(p.Name),
		Email: string(p.Email),
		Phone: string(p.Phone),
	}
}

// DecodePersonRow Turns Storage back into Business.
func DecodePersonRow(r PersonRow) Person {
	return Person{
		ID:    PersonID(r.ID),
		Name:  FullName(r.Name),
		Email: EmailAddress(r.Email),
		Phone: PhoneNumber(r.Phone),
	}
}

// DecodePeopleRows Repeats the Move for a whole Page.
func DecodePeopleRows(rows []PersonRow) []Person {
	people := make([]Person, 0, len(rows))
	for _, row := range rows {
		people = append(people, DecodePersonRow(row))
	}
	return people
}
