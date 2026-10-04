package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/faults"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/school"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
)

// The Form of a Request, Checked before its Meaning.
var (
	ErrBodyIsBroken = faults.RefuseInvalidInput("body is not valid JSON")
	ErrPathIsBroken = faults.RefuseInvalidInput("path Holds no valid Identity")
)

// ReadPathNumber Reads the positive whole Number in the Path's {id}.
func ReadPathNumber(req transport.Request) (uint, error) {
	number, err := strconv.ParseUint(req.Path["id"], 10, 32)
	if err != nil || number == 0 {
		return 0, ErrPathIsBroken
	}

	return uint(number), nil
}

// ReadJSONBody Reads one JSON Object whose Fields are all there, all known
// and all of the right Type. Anything Less is Body is Broken.
func ReadJSONBody[T any](req transport.Request) (T, error) {
	var body T

	if err := checkFieldsArePresent[T](req.Body); err != nil {
		return body, ErrBodyIsBroken
	}

	decoder := json.NewDecoder(bytes.NewReader(req.Body))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&body); err != nil {
		return body, ErrBodyIsBroken
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return body, ErrBodyIsBroken
	}

	return body, nil
}

// checkFieldsArePresent Refuses a Body with a Field missing or null.
func checkFieldsArePresent[T any](data []byte) error {
	var sent map[string]json.RawMessage
	if err := json.Unmarshal(data, &sent); err != nil || sent == nil {
		return ErrBodyIsBroken
	}

	shape := reflect.TypeFor[T]()
	for position := 0; position < shape.NumField(); position++ {
		name, _, _ := strings.Cut(shape.Field(position).Tag.Get("json"), ",")
		value, present := sent[name]
		if !present || string(value) == "null" {
			return ErrBodyIsBroken
		}
	}

	return nil
}

// ReadPageRequest Reads the Window a list Route Asks for.
// Neither Named Asks for the whole Set; one alone takes the other's Default.
func ReadPageRequest(req transport.Request) (school.Page, error) {
	page, hasPage := req.Query["page"]
	size, hasSize := req.Query["size"]
	if !hasPage && !hasSize {
		return school.Page{}, nil
	}

	window := school.Page{Number: 0, Size: school.DefaultPageSize}
	if hasPage {
		number, err := strconv.Atoi(page)
		if err != nil || number < 0 {
			return school.Page{}, school.ErrPageIsInvalid
		}
		window.Number = number
	}

	if hasSize {
		count, err := strconv.Atoi(size)
		if err != nil || count < 1 {
			return school.Page{}, school.ErrPageIsInvalid
		}
		window.Size = count
	}

	return window, nil
}
