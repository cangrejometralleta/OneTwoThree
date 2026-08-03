package main

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// StudentID and its Siblings Name a Row that already Exists.
type StudentID uint
type CourseID uint

// RUT is the Chilean Tax Identity, Digit included.
// A String Underneath, and still Impossible to Pass as a Name.
type RUT string

// FullName is how a Person Wants to be Called.
type FullName string

// CourseCode is the short Label a Course Answers to.
type CourseCode string

// Student is the Business Truth about one Enrolment.
type Student struct {
	ID     StudentID
	Rut    RUT
	Name   FullName
	Age    int
	Course CourseID
}

// Course is the Business Truth about one Class.
type Course struct {
	ID   CourseID
	Code CourseCode
	Name FullName
}

// The Business Fails in named Ways, never in Numbers.
var (
	ErrNameIsEmpty    = errors.New("name is Empty")
	ErrRutIsInvalid   = errors.New("rut Fails its Check Digit")
	ErrRutTaken       = errors.New("rut is already Registered")
	ErrAgeIsTooLow    = errors.New("age Must be eighteen or more")
	ErrStudentUnknown = errors.New("student not Found")
	ErrCourseUnknown  = errors.New("course not Found")
	ErrPageIsInvalid  = errors.New("page Numbers Must not be negative")
	ErrBodyIsBroken   = errors.New("body is not valid JSON")
	ErrPathIsBroken   = errors.New("path Holds no valid Identity")
)

var rutShape = regexp.MustCompile(`^[0-9]+-[0-9kK]$`)

// CheckStudentRecord Refuses a Student the Business cannot Use.
func (s Student) CheckStudentRecord() error {
	if s.Name == "" {
		return ErrNameIsEmpty
	}

	if !s.Rut.LooksValid() {
		return ErrRutIsInvalid
	}

	return s.CheckStudentAge()
}

// CheckStudentAge Holds the one Rule the School will not Bend.
func (s Student) CheckStudentAge() error {
	if s.Age < 18 {
		return ErrAgeIsTooLow
	}

	return nil
}

// LooksValid Runs the Modulo eleven Check the Digit Encodes.
// Reference: https://es.wikipedia.org/wiki/Rol_%C3%9Anico_Tributario
func (r RUT) LooksValid() bool {
	if !rutShape.MatchString(string(r)) {
		return false
	}

	body, digit, _ := strings.Cut(strings.ToLower(string(r)), "-")

	return digit == CheckDigitFor(body)
}

// CheckDigitFor Folds a RUT Body into its single Check Character.
func CheckDigitFor(body string) string {
	number, err := strconv.Atoi(body)
	if err != nil {
		return ""
	}

	sum, index := 0, 0
	for ; number > 0; number /= 10 {
		sum += number % 10 * RutWeights[index%len(RutWeights)]
		index++
	}

	return NameCheckRemainder(11 - sum%11)
}

// RutWeights Cycles two through seven, Read right to left.
var RutWeights = []int{2, 3, 4, 5, 6, 7}

// NameCheckRemainder Turns a Remainder into the Character it Means.
func NameCheckRemainder(remainder int) string {
	switch remainder {
	case 11:
		return "0"
	case 10:
		return "k"
	default:
		return strconv.Itoa(remainder)
	}
}
