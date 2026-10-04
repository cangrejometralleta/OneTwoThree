package serving

import (
	"net/http"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
)

// StdlibServer Serves the Routes with net/http alone.
type StdlibServer struct{}

// BuildHandler Declares each Route on a ServeMux, the Method in its Pattern.
func (StdlibServer) BuildHandler(routes []transport.Route) http.Handler {
	mux := http.NewServeMux()
	for _, route := range routes {
		mux.HandleFunc(route.Method+" "+route.Pattern, ServeRouteHandler(route, func(r *http.Request, name string) string {
			return r.PathValue(name)
		}))
	}

	return mux
}

// ServeRoutes Listens on the Address until it Fails.
func (s StdlibServer) ServeRoutes(routes []transport.Route, address string) error {
	return http.ListenAndServe(address, s.BuildHandler(routes))
}
