package transport

// This Package Holds the Shapes the Layers Speak in.
// It Imports nothing, so every Layer can Depend on it.

// Request is what a Handler Receives.
// No Framework Type Appears here.
type Request struct {
	Path  map[string]string
	Query map[string]string
	Token string
	Body  []byte

	// Caller is Empty until the Application Names it.
	// A Handler Reads it; no Adapter ever Fills it.
	Caller string
}

// Response is what a Handler Returns.
type Response struct {
	Status int
	Body   any
}

// Handler is the only Signature the Business Layer Knows.
type Handler func(Request) Response

// Route Binds one Method and one Pattern to one Handler.
// Patterns Use {name}, and each Adapter Translates from there.
type Route struct {
	Method  string
	Pattern string
	Handle  Handler
}

// Server is the Contract every HTTP Adapter Fulfils.
// Swapping Frameworks Means Swapping the Value behind this Interface.
type Server interface {
	ServeRoutes(routes []Route, address string) error
}
