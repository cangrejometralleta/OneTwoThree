package serving

import (
	"net/http"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
)

// StdlibServer Serves a Libretto with net/http alone.
type StdlibServer struct{}

// BuildHandler Hands every Request to the Libretto; no Mux Decides anything.
func (StdlibServer) BuildHandler(libretto transport.Libretto) http.Handler {
	return ServeLibretto(libretto)
}

// ServeRoutes Listens on the Address until it Fails.
func (s StdlibServer) ServeRoutes(libretto transport.Libretto, address string) error {
	return http.ListenAndServe(address, s.BuildHandler(libretto))
}
