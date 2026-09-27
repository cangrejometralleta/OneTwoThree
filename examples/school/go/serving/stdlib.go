package serving

import (
	"encoding/json"
	"net/http"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
)

// StdlibServer Serves the Routes with net/http alone.
// Reference: https://pkg.go.dev/net/http#ServeMux
type StdlibServer struct{}

// ServeRoutes Mounts every Route on a stdlib Mux.
func (StdlibServer) ServeRoutes(routes []transport.Route, address string) error {
	mux := http.NewServeMux()

	for _, route := range routes {
		mux.HandleFunc(route.Method+" "+route.Pattern, StdlibHandlerFor(route))
	}

	return http.ListenAndServe(address, mux)
}

// StdlibHandlerFor Wraps one Route into a stdlib HandlerFunc.
func StdlibHandlerFor(route transport.Route) http.HandlerFunc {
	names := ListPatternParams(route.Pattern)

	return func(w http.ResponseWriter, r *http.Request) {
		path := map[string]string{}
		for _, name := range names {
			path[name] = r.PathValue(name)
		}

		WriteReplyAsJSON(w, route.Handle(ReadRequestValues(r, path)))
	}
}

// WriteReplyAsJSON is the one Place that Touches a ResponseWriter.
func WriteReplyAsJSON(w http.ResponseWriter, reply transport.Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(reply.Status)

	if reply.Body != nil {
		_ = json.NewEncoder(w).Encode(reply.Body)
	}
}
