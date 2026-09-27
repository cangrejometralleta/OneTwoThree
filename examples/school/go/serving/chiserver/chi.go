package chiserver

// This Package is the only one that Imports chi.
// Swapping it out Edits this Package and main.go, and nothing else.

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/serving"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
)

// ChiServer Serves the same Routes with chi.
// Reference: https://pkg.go.dev/github.com/go-chi/chi/v5
type ChiServer struct{}

// ServeRoutes Mounts every Route on a chi Router.
func (ChiServer) ServeRoutes(routes []transport.Route, address string) error {
	router := chi.NewRouter()

	for _, route := range routes {
		router.MethodFunc(route.Method, route.Pattern, ChiHandlerFor(route))
	}

	return http.ListenAndServe(address, router)
}

// ChiHandlerFor Wraps one Route, Reading Params the chi Way.
func ChiHandlerFor(route transport.Route) http.HandlerFunc {
	names := serving.ListPatternParams(route.Pattern)

	return func(w http.ResponseWriter, r *http.Request) {
		path := map[string]string{}
		for _, name := range names {
			path[name] = chi.URLParam(r, name)
		}

		serving.WriteReplyAsJSON(w, route.Handle(serving.ReadRequestValues(r, path)))
	}
}
