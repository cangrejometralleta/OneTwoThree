package transport

// Request is what the Door Hands the Script: no Framework Type Crosses it.
type Request struct {
	Path   map[string]string
	Query  map[string]string
	Body   []byte
	Token  string
	Caller string
}

// Response is what the Script Hands back: a Status and a Body to Render.
type Response struct {
	Status int
	Body   any
}

// Handler Answers one Request.
type Handler func(Request) Response

// Route Declares one Door: Method, Pattern and the Handler behind it.
type Route struct {
	Method  string
	Pattern string
	Handle  Handler
}

// Server Serves Routes on an Address; each Framework Fills it in its own Package.
type Server interface {
	ServeRoutes(routes []Route, address string) error
}
