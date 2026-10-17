package tokens

import (
	"encoding/base64"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/faults"
)

// The Spec's two Tokens, Signed with spec-secret: one Lives until 2100, one Died at exp 1.
const (
	specLiving = "eyJzdWIiOiJzdHVkZW50LXJlZ2lzdHJ5IiwiZXhwIjo0MTAyNDQ0ODAwfQ.fgl9OtJeU8f90EIkPjELn34yWKjIbn3QcYn5UFoLAyY"
	specDead   = "eyJzdWIiOiJzdHVkZW50LXJlZ2lzdHJ5IiwiZXhwIjoxfQ.V8VQ8Ie_fo88z4EBJ0QbtIc5ONTOqQ_-l_WLuhwz5t4"
)

func clockedAt(secret string, now time.Time) AccessTokens {
	tokens := BuildAccessTokens(secret, 60)
	tokens.now = func() time.Time { return now }

	return tokens
}

var issuedAt = time.Unix(1_000_000, 0)

func TestAccessTokenDiesOnlyAfterItsDeadline(t *testing.T) {
	token, err := clockedAt("secret", issuedAt).IssueAccessToken("student-registry")
	if err != nil {
		t.Fatal(err)
	}

	for _, after := range []time.Duration{0, 59 * time.Second, 60 * time.Second, 60*time.Second + 999*time.Millisecond} {
		if caller, err := clockedAt("secret", issuedAt.Add(after)).ReadCaller(token); err != nil || caller != "student-registry" {
			t.Errorf("%v after: a living Token must Name its Caller, caller=%q err=%v", after, caller, err)
		}
	}

	if _, err := clockedAt("secret", issuedAt.Add(61*time.Second)).ReadCaller(token); !errors.Is(err, ErrTokenIsInvalid) {
		t.Fatalf("a Token past its Deadline must be Dead, err=%v", err)
	}
}

func TestAMintedTokenCarriesCompactClaimsInTheSpecsOrder(t *testing.T) {
	token, _ := clockedAt("secret", issuedAt).IssueAccessToken("student-registry")
	want := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"student-registry","exp":1000060}`))

	if claims, _, _ := strings.Cut(token, "."); claims != want {
		t.Fatalf("claims=%s", claims)
	}
}

func TestTheSpecsTokensReadTheSameInEveryRuntime(t *testing.T) {
	spec := clockedAt("spec-secret", time.Date(2026, 10, 17, 0, 0, 0, 0, time.UTC))

	if caller, err := spec.ReadCaller(specLiving); err != nil || caller != "student-registry" {
		t.Fatalf("caller=%q err=%v", caller, err)
	}
	if _, err := spec.ReadCaller(specDead); !errors.Is(err, ErrTokenIsInvalid) {
		t.Fatalf("err=%v", err)
	}
}

func TestATamperedClaimBreaksTheSignature(t *testing.T) {
	token, _ := clockedAt("secret", issuedAt).IssueAccessToken("student-registry")
	_, signature, _ := strings.Cut(token, ".")
	forged := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"admin","exp":99999999999}`)) + "." + signature

	if _, err := clockedAt("secret", issuedAt).ReadCaller(forged); !errors.Is(err, ErrTokenIsInvalid) {
		t.Fatalf("a changed Claim must Break the Signature, err=%v", err)
	}
}

func TestAnotherSecretAndAnythingUnsignedAreRefused(t *testing.T) {
	token, _ := clockedAt("secret", issuedAt).IssueAccessToken("student-registry")
	other := clockedAt("another", issuedAt)

	for name, candidate := range map[string]string{
		"another secret": token,
		"unsigned":       strings.Split(token, ".")[0],
		"empty":          "",
		"just a dot":     ".",
		"not base64":     "%%%.%%%",
		"three parts":    token + ".x",
	} {
		if _, err := other.ReadCaller(candidate); !errors.Is(err, ErrTokenIsInvalid) {
			t.Errorf("%s must be Refused, err=%v", name, err)
		}
	}
}

func TestATokenWithNoSubjectIsRefused(t *testing.T) {
	tokens := clockedAt("secret", issuedAt)
	encoded := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"","exp":99999999999}`))

	if _, err := tokens.ReadCaller(encoded + "." + tokens.sign(encoded)); !errors.Is(err, ErrTokenIsInvalid) {
		t.Fatalf("a Token that Names no one must be Refused, err=%v", err)
	}
}

func TestAnEmptySecretSignsNothingAndTrustsNothing(t *testing.T) {
	if _, err := clockedAt("", issuedAt).IssueAccessToken("student-registry"); !errors.Is(err, ErrTokenIsInvalid) {
		t.Fatalf("an empty Secret signs nothing, err=%v", err)
	}
	if _, err := clockedAt("secret", issuedAt).IssueAccessToken(""); !errors.Is(err, ErrTokenIsInvalid) {
		t.Fatalf("a Token must Name its Subject, err=%v", err)
	}
	if _, err := clockedAt("", issuedAt).ReadCaller(specLiving); !errors.Is(err, ErrTokenIsInvalid) {
		t.Fatalf("an empty Secret must Trust nothing, err=%v", err)
	}
}

func TestEveryTokenFailureAnswersAsAnUnprovenCaller(t *testing.T) {
	if got := faults.ReadFaultKind(ErrTokenIsInvalid); got != faults.UnprovenCaller {
		t.Fatalf("wanted UnprovenCaller, got %d", got)
	}
}
