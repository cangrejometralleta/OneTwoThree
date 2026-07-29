package main

import "errors"

// PersonID Names a Row that already Exists.
type PersonID uint

// FullName is how a Person Wants to be Called.
type FullName string

// EmailAddress Reaches a Person.
type EmailAddress string

// PhoneNumber Reaches a Person faster.
type PhoneNumber string

// Person is the Business Truth.
// Every Field Carries its own Type, so a Compiler Refuses
// to Pass an Email where a Name Belongs.
type Person struct {
	ID    PersonID
	Name  FullName
	Email EmailAddress
	Phone PhoneNumber
}

// The Business Fails in exactly five Ways, and each one has a Name.
var (
	ErrNameIsEmpty   = errors.New("name is Empty")
	ErrMailIsEmpty   = errors.New("email is Empty")
	ErrPersonUnknown = errors.New("person not Found")
	ErrBodyIsBroken  = errors.New("body is not valid JSON")
	ErrPathIsBroken  = errors.New("path Holds no valid Identity")
)

// CheckPersonRecord Refuses a Person the Business cannot Use.
func (p Person) CheckPersonRecord() error {
	if p.Name == "" {
		return ErrNameIsEmpty
	}

	if p.Email == "" {
		return ErrMailIsEmpty
	}

	return nil
}
