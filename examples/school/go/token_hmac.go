package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"time"
)

// ErrTokenIsInvalid Covers every Way a Token can Fail.
var ErrTokenIsInvalid = errors.New("token is Invalid or Expired")

// HmacTokens Fulfils TokenIssuer with the standard Library only.
// Reference: https://pkg.go.dev/crypto/hmac
type HmacTokens struct {
	Secret []byte
	Life   time.Duration
	Now    func() time.Time
}

// IssueAccessToken Signs a Subject together with its Deadline.
func (h HmacTokens) IssueAccessToken(subject string) (string, error) {
	if len(h.Secret) == 0 {
		return "", ErrTokenIsInvalid
	}

	claim := subject + "." + strconv.FormatInt(h.Now().Add(h.Life).Unix(), 10)

	return claim + "." + h.SignTokenClaim(claim), nil
}

// ReadAccessToken Returns the Subject, or Says why it Cannot.
// The Comparison Reads forward: a Token Dies once its Deadline Passes.
func (h HmacTokens) ReadAccessToken(token string) (string, error) {
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
func (h HmacTokens) OpenSignedClaim(token string) (string, error) {
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
func (h HmacTokens) SignTokenClaim(claim string) string {
	mac := hmac.New(sha256.New, h.Secret)
	mac.Write([]byte(claim))

	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
