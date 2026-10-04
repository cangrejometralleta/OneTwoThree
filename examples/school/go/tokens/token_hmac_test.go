package tokens

import (
	"encoding/base64"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/faults"
)

func buildClockedTokens(secret string, now *time.Time) AccessTokens {
	tokens := BuildAccessTokens(secret, 60)
	tokens.now = func() time.Time { return *now }

	return tokens
}

// A Token Dies once its Deadline Passes, and never before.
func TestAccessTokenDiesOnlyAfterItsDeadline(t *testing.T) {
	now := time.Unix(1_000_000, 0)
	tokens := buildClockedTokens("secret", &now)
	token, err := tokens.IssueAccessToken("student-registry")
	if err != nil {
		t.Fatal(err)
	}

	issuedAt := now
	for _, living := range []struct {
		name  string
		after time.Duration
	}{
		{"just issued", 0},
		{"a second before", 59 * time.Second},
		{"exactly the deadline", 60 * time.Second},
	} {
		now = issuedAt.Add(living.after)
		if caller, err := tokens.ReadCaller(token); err != nil || caller != "student-registry" {
			t.Errorf("%s: a living Token must Name its Caller, caller=%q err=%v", living.name, caller, err)
		}
	}

	now = now.Add(time.Second)
	if _, err := tokens.ReadCaller(token); !errors.Is(err, ErrTokenIsInvalid) {
		t.Fatalf("a Token past its Deadline must be Dead, err=%v", err)
	}
}

func TestAccessTokenRefusesATamperedClaim(t *testing.T) {
	now := time.Unix(1_000_000, 0)
	tokens := buildClockedTokens("secret", &now)
	token, _ := tokens.IssueAccessToken("student-registry")
	_, signature, _ := strings.Cut(token, ".")
	forged := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"admin","exp":99999999999}`)) + "." + signature

	if _, err := tokens.ReadCaller(forged); !errors.Is(err, ErrTokenIsInvalid) {
		t.Fatalf("a changed Claim must Break the Signature, err=%v", err)
	}
}

func TestAccessTokenRefusesAnotherSecretAndNoSignature(t *testing.T) {
	now := time.Unix(1_000_000, 0)
	token, _ := buildClockedTokens("secret", &now).IssueAccessToken("student-registry")
	other := buildClockedTokens("another", &now)

	for name, candidate := range map[string]string{
		"another secret": token,
		"unsigned":       strings.Split(token, ".")[0],
		"empty":          "",
		"just a dot":     ".",
		"not base64":     "%%%.%%%",
	} {
		if _, err := other.ReadCaller(candidate); !errors.Is(err, ErrTokenIsInvalid) {
			t.Errorf("%s must be Refused, err=%v", name, err)
		}
	}
}

func TestAccessTokenRefusesATokenWithNoSubject(t *testing.T) {
	now := time.Unix(1_000_000, 0)
	tokens := buildClockedTokens("secret", &now)
	encoded := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"","exp":99999999999}`))
	signed := encoded + "." + tokens.sign(encoded)

	if _, err := tokens.ReadCaller(signed); !errors.Is(err, ErrTokenIsInvalid) {
		t.Fatalf("a Token that Names no one must be Refused, err=%v", err)
	}
}

func TestEveryTokenFailureAnswersUnauthorised(t *testing.T) {
	if got := faults.ReadFaultKind(ErrTokenIsInvalid); got != faults.UnprovenCaller {
		t.Fatalf("wanted UnprovenCaller, got %d", got)
	}
}

func TestIssueAccessTokenRefusesAnEmptySecretAndNoSubject(t *testing.T) {
	if _, err := BuildAccessTokens("", 60).IssueAccessToken("student-registry"); !errors.Is(err, ErrTokenIsInvalid) {
		t.Fatalf("an empty Secret signs nothing, err=%v", err)
	}
	if _, err := BuildAccessTokens("secret", 60).IssueAccessToken(""); !errors.Is(err, ErrTokenIsInvalid) {
		t.Fatalf("a Token must Name its Subject, err=%v", err)
	}
	if _, err := BuildAccessTokens("", 60).ReadCaller("a.b"); !errors.Is(err, ErrTokenIsInvalid) {
		t.Fatalf("an empty Secret must Trust nothing, err=%v", err)
	}
}
