package app

import (
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/school"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
)

// RequireProvenCaller Names the Caller before the Story Starts.
// A Handler behind this Reads req.Caller and Trusts it.
func RequireProvenCaller(tokens school.TokenIssuer, tell Telling) Telling {
	return func(req transport.Request) (any, error) {
		caller, err := tokens.ReadAccessToken(req.Token)
		if err != nil {
			return nil, err
		}

		req.Caller = caller

		return tell(req)
	}
}
