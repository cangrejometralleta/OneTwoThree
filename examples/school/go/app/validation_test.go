package app

import (
	"bytes"
	"errors"
	"math"
	"testing"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/school"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/wire"
)

const whole = `{"rut":"11111111-1","name":"Ana","age":20,"courseId":1}`

func pathOf(id string) transport.Request { return transport.Request{Path: map[string]string{"id": id}} }

func queryOf(query map[string]string) transport.Request { return transport.Request{Query: query} }

func TestReadPathNumberWantsAPositiveWholeNumber(t *testing.T) {
	for text, want := range map[string]int64{"7": 7, "012": 12} {
		if got, err := ReadPathNumber(pathOf(text)); err != nil || got != want {
			t.Errorf("%s: got=%d err=%v", text, got, err)
		}
	}

	for _, bad := range []string{"", "abc", "0", "000", "-3", "1.5", "+1", " 1", "1e1", "0x1"} {
		if _, err := ReadPathNumber(pathOf(bad)); !errors.Is(err, ErrPathIsBroken) {
			t.Errorf("%q must be Path is Broken, err=%v", bad, err)
		}
	}
}

func TestAPathNumberBeyondSixtyFourBitsNamesNobodyRatherThanBreaking(t *testing.T) {
	if got, err := ReadPathNumber(pathOf("99999999999999999999999")); err != nil || got != math.MaxInt64 {
		t.Fatalf("got=%d err=%v", got, err)
	}
}

func TestReadJSONBodyRefusesAnythingButOneWholeObject(t *testing.T) {
	if body, err := ReadJSONBody[wire.StudentBody](transport.Request{Body: []byte(whole)}); err != nil || body.Course != 1 {
		t.Fatalf("a whole Body must Pass, body=%v err=%v", body, err)
	}

	broken := map[string]string{
		"not json":               `{"rut":`,
		"empty":                  ``,
		"missing field":          `{"rut":"11111111-1","name":"Ana","age":20}`,
		"null field":             `{"rut":"11111111-1","name":null,"age":20,"courseId":1}`,
		"unknown field":          `{"rut":"11111111-1","name":"Ana","age":20,"courseId":1,"id":9}`,
		"wrong type":             `{"rut":"11111111-1","name":"Ana","age":"20","courseId":1}`,
		"text as number":         `{"rut":11111111,"name":"Ana","age":20,"courseId":1}`,
		"fractional":             `{"rut":"11111111-1","name":"Ana","age":20.5,"courseId":1}`,
		"whole written as float": `{"rut":"11111111-1","name":"Ana","age":20.0,"courseId":1}`,
		"beyond 64 bits":         `{"rut":"11111111-1","name":"Ana","age":20,"courseId":9223372036854775808}`,
		"second value":           whole + ` {}`,
		"an array":               `[]`,
		"null document":          `null`,
		"repeated key":           `{"rut":"11111111-1","rut":"x","name":"Ana","age":20,"courseId":1}`,
	}

	for name, body := range broken {
		if _, err := ReadJSONBody[wire.StudentBody](transport.Request{Body: []byte(body)}); !errors.Is(err, ErrBodyIsBroken) {
			t.Errorf("%s must be Body is Broken, err=%v", name, err)
		}
	}
}

func TestANegativeCourseIsWellFormedAndLeftToTheCourseCheck(t *testing.T) {
	body, err := ReadJSONBody[wire.StudentBody](transport.Request{Body: []byte(`{"rut":"11111111-1","name":"Ana","age":20,"courseId":-1}`)})
	if err != nil || body.Course != -1 {
		t.Fatalf("body=%v err=%v", body, err)
	}
}

func TestABodyMustBeUTF8AndAtMostOneMebibyte(t *testing.T) {
	latin := bytes.Replace([]byte(whole), []byte("Ana"), []byte("An\xe1"), 1)
	if _, err := ReadJSONBody[wire.StudentBody](transport.Request{Body: latin}); !errors.Is(err, ErrBodyIsBroken) {
		t.Errorf("latin-1 must be Body is Broken, err=%v", err)
	}

	oversized := append([]byte(whole), bytes.Repeat([]byte(" "), transport.MaxBodyBytes)...)
	if _, err := ReadJSONBody[wire.StudentBody](transport.Request{Body: oversized}); !errors.Is(err, ErrBodyIsBroken) {
		t.Errorf("a Body past the Bound must be Body is Broken, err=%v", err)
	}
}

func TestReadPageRequestWindowsOnlyWhatWasAsked(t *testing.T) {
	cases := map[string]struct {
		query map[string]string
		want  school.Page
	}{
		"neither is the whole set":  {map[string]string{}, school.WholeSet},
		"an unknown key is ignored": {map[string]string{"sort": "name"}, school.WholeSet},
		"both":                      {map[string]string{"page": "2", "size": "5"}, school.Page{Number: 2, Size: 5}},
		"page alone":                {map[string]string{"page": "3"}, school.Page{Number: 3, Size: school.DefaultPageSize}},
		"size alone":                {map[string]string{"size": "25"}, school.Page{Number: 0, Size: 25}},
		"beyond 64 bits":            {map[string]string{"page": "99999999999999999999", "size": "10"}, school.Page{Number: math.MaxInt64, Size: 10}},
	}

	for name, test := range cases {
		if got, err := ReadPageRequest(queryOf(test.query)); err != nil || got != test.want {
			t.Errorf("%s: wanted %v, got %v err=%v", name, test.want, got, err)
		}
	}
}

func TestReadPageRequestRefusesWhatIsNotDigitsInRange(t *testing.T) {
	for _, query := range []map[string]string{
		{"page": "abc"}, {"page": "-1"}, {"page": ""}, {"page": "1.5"}, {"page": "+1"}, {"page": "1e1"},
		{"page": " 1"}, {"page": "0x1"}, {"size": "many"}, {"size": "0"}, {"size": "-4"}, {"size": ""},
	} {
		if _, err := ReadPageRequest(queryOf(query)); !errors.Is(err, school.ErrPageIsInvalid) {
			t.Errorf("%v must be Page is Invalid, err=%v", query, err)
		}
	}
}

func TestABearerTokenIsReadOnlyAfterTheScheme(t *testing.T) {
	cases := map[string]string{"bEaReR abc": "abc", "Basic abc": "", "": ""}

	for header, want := range cases {
		if got := ReadBearerToken(bearing(header)); got != want {
			t.Errorf("%q: wanted %q, got %q", header, want, got)
		}
	}
}
