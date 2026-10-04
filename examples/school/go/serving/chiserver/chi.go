package chiserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/serving"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
)

// ChiServer Serves the Routes with chi. This Package is the only one that Imports chi.
type ChiServer struct{}

// BuildHandler Declares each Route on a chi Router; its Pattern already Reads {id}.
func (ChiServer) BuildHandler(routes []transport.Route) http.Handler {
	router := chi.NewRouter()
	for _, route := range routes {
		router.MethodFunc(route.Method, route.Pattern, serving.ServeRouteHandler(route, func(r *http.Request, name string) string {
			return chi.URLParam(r, name)
		}))
	}

	return router
}

// ServeRoutes Listens on the Address until it Fails.
func (s ChiServer) ServeRoutes(routes []transport.Route, address string) error {
	return http.ListenAndServe(address, s.BuildHandler(routes))
}
