package main

import (
	"io"
	"net/http"
	"strings"
)

// Request is what a Handler Receives.
// No Framework Type Appears here, and that is the whole Point.
type Request struct {
	Path  map[string]string
	Query map[string]string
	Body  []byte
}

// Response is what a Handler Returns.
type Response struct {
	Status int
	Body   any
}

// Handler is the only Signature the Business Layer Knows.
type Handler func(Request) Response

// Route Binds one Method and one Pattern to one Handler.
// Patterns Use {name}, and each Adapter Translates from there.
type Route struct {
	Method  string
	Pattern string
	Handle  Handler
}

// Server is the Contract every Adapter Fulfils.
// Swapping Frameworks Means Swapping the Value behind this Interface.
type Server interface {
	ServeRoutes(routes []Route, address string) error
}

// ReadRequestValues Copies an http Request into our own Shape.
func ReadRequestValues(r *http.Request, path map[string]string) Request {
	body, _ := io.ReadAll(r.Body)

	query := map[string]string{}
	for key, values := range r.URL.Query() {
		query[key] = values[0]
	}

	return Request{Path: path, Query: query, Body: body}
}

// ListPatternParams Finds every {name} inside a Route Pattern.
func ListPatternParams(pattern string) []string {
	found := []string{}

	for _, part := range strings.Split(pattern, "/") {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			found = append(found, strings.Trim(part, "{}"))
		}
	}

	return found
}
