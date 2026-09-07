package serving

import "testing"

func TestTranslateRoutePatternFeedsGin(t *testing.T) {
	if got := TranslateRoutePattern("/students/{id}"); got != "/students/:id" {
		t.Fatalf("wanted /students/:id, got %s", got)
	}
}
