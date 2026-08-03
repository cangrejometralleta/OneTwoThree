package main

import (
	"encoding/json"
	"net/http"
)

// ServeOrderRoutes Mounts every Route on a stdlib Mux.
// The only File in this Program that Imports net/http as a Framework Choice.
// One Service, one Framework: the Portability Point already Lives in School.
func ServeOrderRoutes(routes []Route, address string) error {
	mux := http.NewServeMux()

	for _, route := range routes {
		mux.HandleFunc(route.Method+" "+route.Pattern, StdlibHandlerFor(route))
	}

	return http.ListenAndServe(address, mux)
}

// StdlibHandlerFor Wraps one Route into a stdlib HandlerFunc.
func StdlibHandlerFor(route Route) http.HandlerFunc {
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
func WriteReplyAsJSON(w http.ResponseWriter, reply Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(reply.Status)

	if reply.Body != nil {
		_ = json.NewEncoder(w).Encode(reply.Body)
	}
}
