package tokens

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/faults"
)

// ErrTokenIsInvalid Answers every Token the Guard cannot Trust, whatever Broke.
// A Caller Learns that it Failed and never why.
var ErrTokenIsInvalid = faults.RefuseUnprovenCaller("token is Invalid or Expired")

// AccessTokens Mints and Reads Bearer Tokens signed with HMAC-SHA256,
// using the standard Library alone. A Token is two base64url Parts joined by
// a dot: the Claims, then the Signature over them.
type AccessTokens struct {
	secret []byte
	life   time.Duration
	now    func() time.Time
}

type claims struct {
	Subject  string `json:"sub"`
	Deadline int64  `json:"exp"`
}

// BuildAccessTokens Casts the Issuer with its Secret and how long a Token Lives.
func BuildAccessTokens(secret string, lifeSeconds int) AccessTokens {
	return AccessTokens{secret: []byte(secret), life: time.Duration(lifeSeconds) * time.Second, now: time.Now}
}

// IssueAccessToken Mints a Token for one Subject. An empty Secret Signs nothing.
func (a AccessTokens) IssueAccessToken(subject string) (string, error) {
	if len(a.secret) == 0 || subject == "" {
		return "", ErrTokenIsInvalid
	}

	payload, err := json.Marshal(claims{Subject: subject, Deadline: a.now().Add(a.life).Unix()})
	if err != nil {
		return "", err
	}

	encoded := base64.RawURLEncoding.EncodeToString(payload)

	return encoded + "." + a.sign(encoded), nil
}

// ReadCaller Names the Subject a Token Proves.
// The Comparison Reads forward: a Token is Dead when now is after its Deadline.
func (a AccessTokens) ReadCaller(token string) (string, error) {
	encoded, signature, found := strings.Cut(token, ".")
	if !found || len(a.secret) == 0 || !hmac.Equal([]byte(signature), []byte(a.sign(encoded))) {
		return "", ErrTokenIsInvalid
	}

	payload, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return "", ErrTokenIsInvalid
	}

	var read claims
	if err := json.Unmarshal(payload, &read); err != nil || read.Subject == "" {
		return "", ErrTokenIsInvalid
	}

	if a.now().Unix() > read.Deadline {
		return "", ErrTokenIsInvalid
	}

	return read.Subject, nil
}

func (a AccessTokens) sign(encoded string) string {
	mac := hmac.New(sha256.New, a.secret)
	mac.Write([]byte(encoded))

	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
