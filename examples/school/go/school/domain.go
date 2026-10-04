package school

import (
	"regexp"
	"strconv"
	"strings"
)

type StudentID uint
type CourseID uint

// RUT is the Chilean tax Number a Student is Known by.
type RUT string

// FullName is what a Person is Called.
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

var rutShape = regexp.MustCompile(`^[0-9]+-[0-9kK]$`)

// RutWeights Cycles two through seven, Read right to left.
var RutWeights = []int{2, 3, 4, 5, 6, 7}

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
	if s.Age < ReadSchoolConstants().MinimumAgeYears {
		return ErrAgeIsTooLow
	}

	return nil
}

// CheckCourseRecord Refuses a Course the Business cannot Use.
func (c Course) CheckCourseRecord() error {
	if c.Code == "" {
		return ErrCodeIsEmpty
	}

	if c.Name == "" {
		return ErrNameIsEmpty
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

	sum, position := 0, 0
	for ; number > 0; number /= 10 {
		sum += number % 10 * RutWeights[position%len(RutWeights)]
		position++
	}

	return NameCheckRemainder(11 - sum%11)
}

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
