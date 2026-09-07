package serving

import (
	"io"
	"net/http"
	"strings"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
)

// ReadRequestValues Copies an http Request into our own Shape.
// Reference: https://pkg.go.dev/net/http#Request
func ReadRequestValues(r *http.Request, path map[string]string) transport.Request {
	body, _ := io.ReadAll(r.Body)

	query := map[string]string{}
	for key, values := range r.URL.Query() {
		query[key] = values[0]
	}

	return transport.Request{Path: path, Query: query, Token: ReadBearerToken(r), Body: body}
}

// ReadBearerToken Pulls the Token out of the Authorization Header.
// Reference: https://datatracker.ietf.org/doc/html/rfc6750#section-2.1
func ReadBearerToken(r *http.Request) string {
	header := r.Header.Get("Authorization")

	return strings.TrimPrefix(header, "Bearer ")
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
