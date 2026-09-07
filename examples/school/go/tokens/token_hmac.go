package tokens

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"strings"
	"time"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/faults"
)

// ErrTokenIsInvalid Covers every Way a Token can Fail.
// A Caller who Cannot Prove itself Gets one Answer, never a Reason why.
var ErrTokenIsInvalid = faults.RefuseUnprovenCaller("token is Invalid or Expired")

// AccessTokens Fulfils TokenIssuer with the standard Library only.
// Reference: https://pkg.go.dev/crypto/hmac
type AccessTokens struct {
	Secret []byte
	Life   time.Duration
	Now    func() time.Time
}

// BuildAccessTokens Takes the two Values Validated at Startup.
// It Names no Config Shape.
func BuildAccessTokens(secret string, lifeSeconds int) AccessTokens {
	return AccessTokens{
		Secret: []byte(secret),
		Life:   time.Duration(lifeSeconds) * time.Second,
		Now:    time.Now,
	}
}

// IssueAccessToken Signs a Subject together with its Deadline.
func (h AccessTokens) IssueAccessToken(subject string) (string, error) {
	if len(h.Secret) == 0 {
		return "", ErrTokenIsInvalid
	}

	claim := subject + "." + strconv.FormatInt(h.Now().Add(h.Life).Unix(), 10)

	return claim + "." + h.SignTokenClaim(claim), nil
}

// ReadAccessToken Returns the Subject, or Says why it Cannot.
// The Comparison Reads forward: a Token Dies once its Deadline Passes.
func (h AccessTokens) ReadAccessToken(token string) (string, error) {
	claim, err := h.OpenSignedClaim(token)
	if err != nil {
		return "", err
	}

	subject, deadline, _ := strings.Cut(claim, ".")
	seconds, parseErr := strconv.ParseInt(deadline, 10, 64)
	if parseErr != nil || h.Now().After(time.Unix(seconds, 0)) {
		return "", ErrTokenIsInvalid
	}

	return subject, nil
}

// OpenSignedClaim Returns the Claim only when the Signature Holds.
func (h AccessTokens) OpenSignedClaim(token string) (string, error) {
	cut := strings.LastIndex(token, ".")
	if cut < 0 {
		return "", ErrTokenIsInvalid
	}

	claim, signature := token[:cut], token[cut+1:]

	if !hmac.Equal([]byte(signature), []byte(h.SignTokenClaim(claim))) {
		return "", ErrTokenIsInvalid
	}

	return claim, nil
}

// SignTokenClaim Reduces a Claim to its Signature.
func (h AccessTokens) SignTokenClaim(claim string) string {
	mac := hmac.New(sha256.New, h.Secret)
	mac.Write([]byte(claim))

	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
