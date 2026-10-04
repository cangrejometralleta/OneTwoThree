package app

import (
	"errors"
	"testing"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/school"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/wire"
)

func TestReadPathNumberWantsAPositiveWholeNumber(t *testing.T) {
	good, err := ReadPathNumber(transport.Request{Path: map[string]string{"id": "7"}})
	if err != nil || good != 7 {
		t.Fatalf("good=%d err=%v", good, err)
	}

	for _, bad := range []string{"", "abc", "0", "-3", "1.5", "99999999999"} {
		if _, err := ReadPathNumber(transport.Request{Path: map[string]string{"id": bad}}); !errors.Is(err, ErrPathIsBroken) {
			t.Errorf("%q must be Path is Broken, err=%v", bad, err)
		}
	}
}

func TestReadJSONBodyRefusesAnythingButOneWholeObject(t *testing.T) {
	whole := `{"rut":"11111111-1","name":"Ana","age":20,"courseId":1}`
	if _, err := ReadJSONBody[wire.StudentBody](transport.Request{Body: []byte(whole)}); err != nil {
		t.Fatalf("a whole Body must Pass, err=%v", err)
	}

	broken := map[string]string{
		"not json":      `{"rut":`,
		"empty":         ``,
		"missing field": `{"rut":"11111111-1","name":"Ana","age":20}`,
		"null field":    `{"rut":"11111111-1","name":null,"age":20,"courseId":1}`,
		"unknown field": `{"rut":"11111111-1","name":"Ana","age":20,"courseId":1,"id":9}`,
		"wrong type":    `{"rut":"11111111-1","name":"Ana","age":"20","courseId":1}`,
		"fractional":    `{"rut":"11111111-1","name":"Ana","age":20.5,"courseId":1}`,
		"second value":  whole + ` {}`,
		"an array":      `[]`,
	}

	for name, body := range broken {
		if _, err := ReadJSONBody[wire.StudentBody](transport.Request{Body: []byte(body)}); !errors.Is(err, ErrBodyIsBroken) {
			t.Errorf("%s must be Body is Broken, err=%v", name, err)
		}
	}
}

func TestReadPageRequestWindowsOnlyWhatWasAsked(t *testing.T) {
	cases := map[string]struct {
		query map[string]string
		want  school.Page
	}{
		"neither is the whole set": {map[string]string{}, school.Page{}},
		"both":                     {map[string]string{"page": "2", "size": "5"}, school.Page{Number: 2, Size: 5}},
		"page alone":               {map[string]string{"page": "3"}, school.Page{Number: 3, Size: school.DefaultPageSize}},
		"size alone":               {map[string]string{"size": "25"}, school.Page{Number: 0, Size: 25}},
	}

	for name, test := range cases {
		got, err := ReadPageRequest(transport.Request{Query: test.query})
		if err != nil || got != test.want {
			t.Errorf("%s: wanted %v, got %v err=%v", name, test.want, got, err)
		}
	}
}

func TestReadPageRequestRefusesWhatIsNotAWholeNumberInRange(t *testing.T) {
	for _, query := range []map[string]string{
		{"page": "abc"}, {"page": "-1"}, {"page": ""}, {"page": "1.5"},
		{"size": "many"}, {"size": "0"}, {"size": "-4"}, {"size": ""},
	} {
		if _, err := ReadPageRequest(transport.Request{Query: query}); !errors.Is(err, school.ErrPageIsInvalid) {
			t.Errorf("%v must be Page is Invalid, err=%v", query, err)
		}
	}
}
