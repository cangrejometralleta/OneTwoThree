package school

import (
	"regexp"
	"strconv"
	"strings"
)

// StudentID is the Identity the Store Gives a Student; the Body never Carries it.
type StudentID int64

// CourseID is the Identity the Store Gives a Course.
type CourseID int64

// RUT is the Chilean tax Number a Student is Known by.
type RUT string

// FullName is what a Person or a Course is Called, Stored as it Arrived.
type FullName string

// CourseCode is the short Label a Course Answers to. Codes are not Unique.
type CourseCode string

// Student is the Business Truth about one Enrolment.
type Student struct {
	ID     StudentID
	Rut    RUT
	Name   FullName
	Age    int64
	Course CourseID
}

// Course is the Business Truth about one Class.
type Course struct {
	ID   CourseID
	Code CourseCode
	Name FullName
}

// SchoolRules Holds the Global Constants the Business will not Bend. No Deployment Overrides them.
type SchoolRules struct {
	MinimumAgeYears int64
}

// RutWeights Cycle two through seven, Read right to left.
var RutWeights = []int{2, 3, 4, 5, 6, 7}

// CheckModulus is the Modulo the Check Digit Encodes.
const CheckModulus = 11

var rutShape = regexp.MustCompile(`^[0-9]+-[0-9kK]$`)

// ReadRUT Promotes Text into a RUT, in lower case: 12345-K and 12345-k are one RUT.
// Leading zeros are not Merged.
func ReadRUT(text string) RUT {
	return RUT(strings.ToLower(text))
}

// CheckStudentRecord Refuses a Student the Business cannot Use: Name, then RUT, then Age.
func (s Student) CheckStudentRecord(rules SchoolRules) error {
	if s.Name.IsEmpty() {
		return ErrNameIsEmpty
	}

	if !s.Rut.LooksValid() {
		return ErrRutIsInvalid
	}

	if s.Age < rules.MinimumAgeYears {
		return ErrAgeIsTooLow(rules)
	}

	return nil
}

// CheckCourseRecord Refuses a Course the Business cannot Use: Code, then Name.
func (c Course) CheckCourseRecord() error {
	if strings.TrimSpace(string(c.Code)) == "" {
		return ErrCodeIsEmpty
	}

	if c.Name.IsEmpty() {
		return ErrNameIsEmpty
	}

	return nil
}

// IsEmpty Holds for a Name of only whitespace; the Text itself is never Trimmed.
func (n FullName) IsEmpty() bool {
	return strings.TrimSpace(string(n)) == ""
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

// CheckDigitFor Folds a RUT Body into its single Check Character, Digit by Digit.
func CheckDigitFor(body string) string {
	sum := 0
	for position := 0; position < len(body); position++ {
		digit := int(body[len(body)-1-position] - '0')
		sum = (sum + digit*RutWeights[position%len(RutWeights)]) % CheckModulus
	}

	return NameCheckRemainder(CheckModulus - sum)
}

// NameCheckRemainder Turns a Remainder into the Character it Means: eleven is 0, ten is k.
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
