package ginserver

// This Package is the only one that Imports gin.
// Swapping it out Edits this Package and main.go, and nothing else.

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/serving"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
)

// GinServer Serves the same Routes with gin.
// Gin Owns its own Context, so this Adapter Carries the most Work.
// Reference: https://pkg.go.dev/github.com/gin-gonic/gin
type GinServer struct{}

// ServeRoutes Mounts every Route on a gin Engine.
func (GinServer) ServeRoutes(routes []transport.Route, address string) error {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	for _, route := range routes {
		engine.Handle(route.Method, TranslateRoutePattern(route.Pattern), GinHandlerFor(route))
	}

	return engine.Run(address)
}

// GinHandlerFor Wraps one Route, Reading Params the gin Way.
func GinHandlerFor(route transport.Route) gin.HandlerFunc {
	names := serving.ListPatternParams(route.Pattern)

	return func(c *gin.Context) {
		path := map[string]string{}
		for _, name := range names {
			path[name] = c.Param(name)
		}

		serving.WriteReplyAsJSON(c.Writer, route.Handle(serving.ReadRequestValues(c.Request, path)))
	}
}

// TranslateRoutePattern Turns {id} into :id, which is all gin Wants.
func TranslateRoutePattern(pattern string) string {
	parts := strings.Split(pattern, "/")

	for index, part := range parts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			parts[index] = ":" + strings.Trim(part, "{}")
		}
	}

	return strings.Join(parts, "/")
}
