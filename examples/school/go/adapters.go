package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
)

// StdlibServer Serves the Routes with net/http alone.
// Reference: https://pkg.go.dev/net/http#ServeMux
type StdlibServer struct{}

// ServeRoutes Mounts every Route on a stdlib Mux.
func (StdlibServer) ServeRoutes(routes []Route, address string) error {
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

// ChiServer Serves the same Routes with chi.
// Reference: https://pkg.go.dev/github.com/go-chi/chi/v5
type ChiServer struct{}

// ServeRoutes Mounts every Route on a chi Router.
func (ChiServer) ServeRoutes(routes []Route, address string) error {
	router := chi.NewRouter()

	for _, route := range routes {
		router.MethodFunc(route.Method, route.Pattern, ChiHandlerFor(route))
	}

	return http.ListenAndServe(address, router)
}

// ChiHandlerFor Wraps one Route, Reading Params the chi Way.
func ChiHandlerFor(route Route) http.HandlerFunc {
	names := ListPatternParams(route.Pattern)

	return func(w http.ResponseWriter, r *http.Request) {
		path := map[string]string{}
		for _, name := range names {
			path[name] = chi.URLParam(r, name)
		}

		WriteReplyAsJSON(w, route.Handle(ReadRequestValues(r, path)))
	}
}

// GinServer Serves the same Routes with gin.
// Gin Owns its own Context, so this Adapter Carries the most Work.
// Reference: https://pkg.go.dev/github.com/gin-gonic/gin
type GinServer struct{}

// ServeRoutes Mounts every Route on a gin Engine.
func (GinServer) ServeRoutes(routes []Route, address string) error {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	for _, route := range routes {
		engine.Handle(route.Method, TranslateRoutePattern(route.Pattern), GinHandlerFor(route))
	}

	return engine.Run(address)
}

// GinHandlerFor Wraps one Route, Reading Params the gin Way.
func GinHandlerFor(route Route) gin.HandlerFunc {
	names := ListPatternParams(route.Pattern)

	return func(c *gin.Context) {
		path := map[string]string{}
		for _, name := range names {
			path[name] = c.Param(name)
		}

		WriteReplyAsJSON(c.Writer, route.Handle(ReadRequestValues(c.Request, path)))
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

// WriteReplyAsJSON is the one Place that Touches a ResponseWriter.
func WriteReplyAsJSON(w http.ResponseWriter, reply Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(reply.Status)

	if reply.Body != nil {
		_ = json.NewEncoder(w).Encode(reply.Body)
	}
}

// SelectServerAdapter Picks the Framework at Startup, never at Compile Time.
func SelectServerAdapter(name string) Server {
	switch name {
	case "gin":
		return GinServer{}
	case "chi":
		return ChiServer{}
	default:
		return StdlibServer{}
	}
}
