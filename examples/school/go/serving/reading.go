package serving

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
)

// MaxBodyBytes Bounds what a Client may Send before the Script Reads it.
const MaxBodyBytes = 1 << 20

var patternNames = regexp.MustCompile(`\{([a-z_]+)\}`)

// ReadPatternNames Lists the Names a Route Pattern Declares: /students/{id} Holds id.
func ReadPatternNames(pattern string) []string {
	var names []string
	for _, found := range patternNames.FindAllStringSubmatch(pattern, -1) {
		names = append(names, found[1])
	}

	return names
}

// ReadRequest Turns the Framework's Request into the Door's own: no Framework Type Crosses.
func ReadRequest(r *http.Request, path map[string]string) transport.Request {
	body, _ := io.ReadAll(io.LimitReader(r.Body, MaxBodyBytes))

	query := map[string]string{}
	for name, values := range r.URL.Query() {
		query[name] = values[0]
	}

	return transport.Request{
		Path: path, Query: query, Body: body,
		Token: strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "),
	}
}

// WriteResponse Renders the Script's Response: JSON, and no Body for a 204.
func WriteResponse(w http.ResponseWriter, reply transport.Response) {
	if reply.Status == http.StatusNoContent || reply.Body == nil {
		w.WriteHeader(reply.Status)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(reply.Status)
	_ = json.NewEncoder(w).Encode(reply.Body)
}

// ServeRouteHandler Joins the three Beats every Adapter Shares: read, answer, write.
func ServeRouteHandler(route transport.Route, readPath func(*http.Request, string) string) http.HandlerFunc {
	names := ReadPatternNames(route.Pattern)

	return func(w http.ResponseWriter, r *http.Request) {
		path := make(map[string]string, len(names))
		for _, name := range names {
			path[name] = readPath(r, name)
		}

		WriteResponse(w, route.Handle(ReadRequest(r, path)))
	}
}
