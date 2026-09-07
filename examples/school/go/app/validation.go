package app

import (
	"encoding/json"
	"strconv"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/faults"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/school"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
)

// A Request can be Malformed without Breaking a single Business Rule.
// Those Failures Belong to the Application.
var (
	ErrBodyIsBroken = faults.RefuseInvalidInput("body is not valid JSON")
	ErrPathIsBroken = faults.RefuseInvalidInput("path Holds no valid Identity")
)

// ReadPathNumber Pulls an Identity out of the Path.
func ReadPathNumber(req transport.Request) (uint64, error) {
	raw, err := strconv.ParseUint(req.Path["id"], 10, 64)
	if err != nil {
		return 0, ErrPathIsBroken
	}

	return raw, nil
}

// ReadJSONBody Decodes the Wire into whatever Shape the Caller Expects.
// It Validates the Form, never the Meaning: the Core Owns the Rules.
func ReadJSONBody[T any](req transport.Request) (T, error) {
	var body T

	if err := json.Unmarshal(req.Body, &body); err != nil {
		return body, ErrBodyIsBroken
	}

	return body, nil
}

// ReadPageRequest Reads Pagination, Defaulting to the whole Set.
func ReadPageRequest(req transport.Request) (school.Page, error) {
	page := school.Page{
		Number: ReadWholeOrRefuse(req.Query["page"]),
		Size:   ReadWholeOrRefuse(req.Query["size"]),
	}

	return page, page.CheckPageBounds()
}

// ReadWholeOrRefuse Reads a whole Number, or Returns one the Bounds Refuse.
// An absent Value Means the whole Set; a Word Means the Caller Mistyped.
// Swallowing the Parse Error would Answer two hundred to Nonsense.
func ReadWholeOrRefuse(value string) int {
	if value == "" {
		return 0
	}

	number, err := strconv.Atoi(value)
	if err != nil {
		return -1
	}

	return number
}
