package ginserver

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/serving"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
)

// GinServer Serves the Routes with gin. This Package is the only one that Imports gin.
type GinServer struct{}

var braceName = regexp.MustCompile(`\{([a-z_]+)\}`)

// BuildHandler Declares each Route on a gin Engine, turning {id} into :id.
func (GinServer) BuildHandler(routes []transport.Route) http.Handler {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	for _, route := range routes {
		handle := serving.ServeRouteHandler(route, func(r *http.Request, name string) string {
			return r.PathValue(name)
		})
		engine.Handle(route.Method, braceName.ReplaceAllString(route.Pattern, ":$1"), func(c *gin.Context) {
			for _, param := range c.Params {
				c.Request.SetPathValue(param.Key, param.Value)
			}
			handle(c.Writer, c.Request)
		})
	}

	return engine
}

// ServeRoutes Listens on the Address until it Fails.
func (s GinServer) ServeRoutes(routes []transport.Route, address string) error {
	return http.ListenAndServe(address, s.BuildHandler(routes))
}
