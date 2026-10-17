package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math"
	"math/big"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/faults"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/school"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
)

// The Form of a Request, Checked before its Meaning.
var (
	ErrBodyIsBroken = faults.RefuseInvalidInput("body is not valid JSON")
	ErrPathIsBroken = faults.RefuseInvalidInput("path Holds no valid Identity")
	ErrRouteUnknown = faults.ReportMissingRecord("route not Found")
)

var (
	digitsOnly  = regexp.MustCompile(`^[0-9]+$`)
	wholeNumber = regexp.MustCompile(`^-?(0|[1-9][0-9]*)$`)
	largest     = big.NewInt(math.MaxInt64)
)

const bearerScheme = "bearer "

// ReadPathNumber Reads the positive whole Number in the Path's {id}; 012 is 12.
// A Number beyond 64 bits Names nobody, so it Reads as the largest, which no one Holds.
func ReadPathNumber(req transport.Request) (int64, error) {
	text := req.Path["id"]
	if !digitsOnly.MatchString(text) {
		return 0, ErrPathIsBroken
	}

	number := readClampedNumber(text)
	if number == 0 {
		return 0, ErrPathIsBroken
	}

	return number, nil
}

// ReadJSONBody Reads one JSON Object, in UTF-8, whose Fields are exactly those of T,
// none null, each a String or a whole Number of 64 bits. Anything Less, or More, is Body is Broken.
func ReadJSONBody[T any](req transport.Request) (T, error) {
	var body T
	if len(req.Body) > transport.MaxBodyBytes || !utf8.Valid(req.Body) {
		return body, ErrBodyIsBroken
	}

	sent, err := readObjectFields(req.Body)
	if err != nil {
		return body, ErrBodyIsBroken
	}

	shape := reflect.ValueOf(&body).Elem()
	if shape.NumField() != len(sent) {
		return body, ErrBodyIsBroken
	}

	for position := 0; position < shape.NumField(); position++ {
		name, _, _ := strings.Cut(shape.Type().Field(position).Tag.Get("json"), ",")
		raw, present := sent[name]
		if !present || !setFieldValue(shape.Field(position), raw) {
			return body, ErrBodyIsBroken
		}
	}

	return body, nil
}

// readObjectFields Walks one Object Key by Key, so a repeated Key or a second Value Breaks it.
func readObjectFields(data []byte) (map[string]json.RawMessage, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	if start, err := decoder.Token(); err != nil || start != json.Delim('{') {
		return nil, ErrBodyIsBroken
	}

	fields := map[string]json.RawMessage{}
	for decoder.More() {
		key, err := decoder.Token()
		name, isText := key.(string)
		if err != nil || !isText {
			return nil, ErrBodyIsBroken
		}

		var raw json.RawMessage
		if _, repeated := fields[name]; repeated || decoder.Decode(&raw) != nil {
			return nil, ErrBodyIsBroken
		}
		fields[name] = raw
	}

	if end, err := decoder.Token(); err != nil || end != json.Delim('}') {
		return nil, ErrBodyIsBroken
	}
	if _, err := decoder.Token(); !errors.Is(err, io.EOF) {
		return nil, ErrBodyIsBroken
	}

	return fields, nil
}

// setFieldValue Accepts a String for a String Field and a whole Number for an int64 Field.
func setFieldValue(field reflect.Value, raw json.RawMessage) bool {
	switch field.Kind() {
	case reflect.String:
		var text string
		if raw[0] != '"' || json.Unmarshal(raw, &text) != nil {
			return false
		}
		field.SetString(text)
	case reflect.Int64:
		if !wholeNumber.Match(raw) {
			return false
		}
		number, err := strconv.ParseInt(string(raw), 10, 64)
		if err != nil {
			return false
		}
		field.SetInt(number)
	default:
		return false
	}

	return true
}

// ReadPageRequest Reads the Window a list Route Asks for.
// Neither Named Asks for the whole Set; one alone Takes the other's Default.
func ReadPageRequest(req transport.Request) (school.Page, error) {
	page, hasPage := req.Query["page"]
	size, hasSize := req.Query["size"]
	if !hasPage && !hasSize {
		return school.WholeSet, nil
	}

	number, count := int64(0), int64(school.DefaultPageSize)
	if hasPage {
		if !digitsOnly.MatchString(page) {
			return school.Page{}, school.ErrPageIsInvalid
		}
		number = readClampedNumber(page)
	}
	if hasSize {
		if !digitsOnly.MatchString(size) {
			return school.Page{}, school.ErrPageIsInvalid
		}
		count = readClampedNumber(size)
	}

	return school.BuildPage(number, count)
}

// ReadBearerToken Reads the Token after a Bearer Scheme of any Case, or Nothing.
func ReadBearerToken(req transport.Request) string {
	header := req.Headers["authorization"]
	if len(header) < len(bearerScheme) || !strings.EqualFold(header[:len(bearerScheme)], bearerScheme) {
		return ""
	}

	return header[len(bearerScheme):]
}

// readClampedNumber Reads Digits as far as 64 bits can Hold them.
func readClampedNumber(digits string) int64 {
	number, _ := new(big.Int).SetString(digits, 10)
	if number.Cmp(largest) > 0 {
		return math.MaxInt64
	}

	return number.Int64()
}
