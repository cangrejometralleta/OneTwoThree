package serving

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
)

// maxSwallowBytes Bounds what is Read and Dropped past the Body, so a Client can Read the Answer.
const maxSwallowBytes = 2 << 20

// ServeLibretto Joins the three Beats every Adapter Shares: read, answer, write.
// Every Framework Hands it the whole Request, so none of them Decides.
func ServeLibretto(libretto transport.Libretto) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(io.LimitReader(r.Body, transport.MaxBodyBytes+1))
		_, _ = io.CopyN(io.Discard, r.Body, maxSwallowBytes)

		headers := map[string]string{}
		for name, values := range r.Header {
			headers[strings.ToLower(name)] = values[0]
		}

		req := transport.Request{Query: ReadQueryValues(r.URL.RawQuery), Headers: headers, Body: body}
		WriteResponse(w, AnswerRequest(libretto, r.Method, r.URL.EscapedPath(), req))
	}
}

// AnswerRequest Finds the Route a Method and a raw Path Name, or Answers the Unknown.
func AnswerRequest(libretto transport.Libretto, method, rawPath string, req transport.Request) transport.Response {
	for _, route := range libretto.Routes {
		if names, matched := MatchRoutePattern(route.Pattern, rawPath); matched && route.Method == method {
			req.Path = names

			return route.Handle(req)
		}
	}

	return libretto.Unknown(req)
}

// MatchRoutePattern Reads the Names a Pattern Declares: /students/{id} Holds id. /students/ is not /students.
func MatchRoutePattern(pattern, rawPath string) (map[string]string, bool) {
	wanted := strings.Split(pattern, "/")
	given := strings.Split(rawPath, "/")
	if len(wanted) != len(given) {
		return nil, false
	}

	names := map[string]string{}
	for position, part := range wanted {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			value, err := url.PathUnescape(given[position])
			if given[position] == "" || err != nil {
				return nil, false
			}
			names[part[1:len(part)-1]] = value
		} else if part != given[position] {
			return nil, false
		}
	}

	return names, true
}

// ReadQueryValues Takes the first Value of a repeated Key; a Pair that cannot be Decoded is Dropped.
func ReadQueryValues(rawQuery string) map[string]string {
	values := map[string]string{}
	if rawQuery == "" {
		return values
	}

	for _, pair := range strings.Split(rawQuery, "&") {
		rawName, rawValue, _ := strings.Cut(pair, "=")
		name, nameErr := url.QueryUnescape(rawName)
		value, valueErr := url.QueryUnescape(rawValue)
		if nameErr != nil || valueErr != nil {
			continue
		}
		if _, taken := values[name]; !taken {
			values[name] = value
		}
	}

	return values
}

// WriteResponse Renders the Script's Response: JSON, and no Body for a 204.
func WriteResponse(w http.ResponseWriter, reply transport.Response) {
	if reply.Status == http.StatusNoContent || reply.Body == nil {
		w.WriteHeader(reply.Status)

		return
	}

	data, err := json.Marshal(reply.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(reply.Status)
	_, _ = w.Write(data)
}
