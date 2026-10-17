package ginserver

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/serving"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
)

// GinServer Serves a Libretto through gin. This Package is the only one that Imports gin.
// Reference: https://gin-gonic.com/en/docs/examples/
type GinServer struct{}

// BuildHandler Declares no Route, so gin Hands every Request to the Libretto.
func (GinServer) BuildHandler(libretto transport.Libretto) http.Handler {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.NoRoute(gin.WrapH(serving.ServeLibretto(libretto)))

	return engine
}

// ServeRoutes Listens on the Address until it Fails.
func (s GinServer) ServeRoutes(libretto transport.Libretto, address string) error {
	return http.ListenAndServe(address, s.BuildHandler(libretto))
}
