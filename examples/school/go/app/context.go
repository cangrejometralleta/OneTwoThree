package app

import "github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"

// CallerReader Reads who a Bearer Token Belongs to.
// The Guard Calls it, so the Guard Declares it.
// Reference: https://pkg.go.dev/crypto/hmac
type CallerReader interface {
	ReadAccessToken(token string) (string, error)
}

// RequireProvenCaller Names the Caller before the Story Starts.
// A Handler behind this Reads req.Caller and Trusts it.
func RequireProvenCaller(tokens CallerReader, tell Telling) Telling {
	return func(req transport.Request) (any, error) {
		caller, err := tokens.ReadAccessToken(req.Token)
		if err != nil {
			return nil, err
		}

		req.Caller = caller

		return tell(req)
	}
}
