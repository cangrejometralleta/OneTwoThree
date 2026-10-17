package chiserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/serving"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
)

// ChiServer Serves a Libretto through chi. This Package is the only one that Imports chi.
// Reference: https://go-chi.io/#/pages/routing
type ChiServer struct{}

// BuildHandler Mounts one catch-all Route, so chi Hands every Request to the Libretto.
func (ChiServer) BuildHandler(libretto transport.Libretto) http.Handler {
	router := chi.NewRouter()
	router.Handle("/*", serving.ServeLibretto(libretto))

	return router
}

// ServeRoutes Listens on the Address until it Fails.
func (s ChiServer) ServeRoutes(libretto transport.Libretto, address string) error {
	return http.ListenAndServe(address, s.BuildHandler(libretto))
}
