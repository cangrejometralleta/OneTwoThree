package app

import "github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"

// CallerReader Names the Caller a Token Proves, or Fails.
type CallerReader interface {
	ReadCaller(token string) (string, error)
}

// RequireProvenCaller Names the Caller before the Story starts.
// A Script behind it Reads req.Caller and Trusts it.
func RequireProvenCaller(reader CallerReader, tell Telling) Telling {
	return func(req transport.Request) (any, error) {
		caller, err := reader.ReadCaller(req.Token)
		if err != nil {
			return nil, err
		}

		req.Caller = caller

		return tell(req)
	}
}
