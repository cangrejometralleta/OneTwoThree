package transport

// MaxBodyBytes Bounds what a Client may Send: one mebibyte.
const MaxBodyBytes = 1 << 20

// Request is what the Door Hands the Script: no Framework Type Crosses it.
// Header Names Arrive in lower case; the Body is at most one Byte past the Bound.
type Request struct {
	Path    map[string]string
	Query   map[string]string
	Headers map[string]string
	Body    []byte
	Caller  string
}

// Response is what the Script Hands back: a Status, and a Body to Render or nil.
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

// Libretto is every Route a Service Declares, and the one Answer for a Door it does not Hold.
// The Unknown Answers before any Guard, because there is no Route to Guard.
type Libretto struct {
	Routes  []Route
	Unknown Handler
}

// Server Serves a Libretto on an Address; each Framework Fills it in its own Package.
type Server interface {
	ServeRoutes(libretto Libretto, address string) error
}
