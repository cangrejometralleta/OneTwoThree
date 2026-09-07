package tokens

import (
	"net/http"
	"testing"
	"time"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/faults"
)

// The Comparison that Broke in the Spring Version, Pinned here.
func TestAccessTokenDiesOnlyAfterItsDeadline(t *testing.T) {
	clock := time.Now()
	tokens := AccessTokens{Secret: []byte("s"), Life: time.Minute, Now: func() time.Time { return clock }}

	token, _ := tokens.IssueAccessToken("registry")

	if _, err := tokens.ReadAccessToken(token); err != nil {
		t.Fatalf("a fresh Token Must be accepted, got %v", err)
	}

	clock = clock.Add(2 * time.Minute)
	if _, err := tokens.ReadAccessToken(token); err == nil {
		t.Fatal("an expired Token Must be Refused")
	}
}

func TestAccessTokenRefusesATamperedClaim(t *testing.T) {
	tokens := AccessTokens{Secret: []byte("s"), Life: time.Minute, Now: time.Now}

	token, _ := tokens.IssueAccessToken("registry")

	if _, err := tokens.ReadAccessToken("admin" + token[8:]); err == nil {
		t.Fatal("a rewritten Subject Must be Refused")
	}
}

// Every Way a Token Fails Arrives as the same Answer.
// A Caller who cannot Prove itself Learns nothing else.
func TestEveryTokenFailureAnswersUnauthorised(t *testing.T) {
	signed, _ := BuildAccessTokens("s", 60).IssueAccessToken("registry")

	cases := map[string]string{
		"someone Sends nothing at all":         "",
		"someone Sends a Word, not a Token":    "hello",
		"someone Sends a Claim with no Sign":   "registry.99999999999",
		"someone Rewrites the Subject":         "admin" + signed[8:],
		"someone Keeps a Token past its Death": signed,
	}

	for story, token := range cases {
		t.Run(story, func(t *testing.T) {
			later := BuildAccessTokens("s", 60)
			later.Now = func() time.Time { return time.Now().Add(time.Hour) }

			_, err := later.ReadAccessToken(token)

			CheckUnprovenCaller(t, err)
		})
	}
}

// A Secret that never Arrived Refuses to Sign at all.
func TestIssueAccessTokenRefusesAnEmptySecret(t *testing.T) {
	_, err := BuildAccessTokens("", 60).IssueAccessToken("registry")

	CheckUnprovenCaller(t, err)
}

// CheckUnprovenCaller Reads the Fault the Adapter Declared.
func CheckUnprovenCaller(t *testing.T, err error) {
	t.Helper()

	if err != ErrTokenIsInvalid {
		t.Fatalf("wanted %v, got %v", ErrTokenIsInvalid, err)
	}

	if got := faults.ReadFaultStatus(err); got != http.StatusUnauthorized {
		t.Fatalf("wanted 401, got %d", got)
	}
}
