package serving

import (
	"testing"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/serving/adaptertest"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
)

var _ transport.Server = StdlibServer{}

func TestStdlibServerAnswersTheSharedContract(t *testing.T) {
	adaptertest.CheckAdapterAnswersTheSameWay(t, StdlibServer{}.BuildHandler)
}

func TestARoutePatternMatchesOnlyItsOwnShape(t *testing.T) {
	if names, matched := MatchRoutePattern("/students/{id}", "/students/%31"); !matched || names["id"] != "1" {
		t.Fatalf("names=%v matched=%v", names, matched)
	}
	for _, path := range []string{"/students/", "/students", "/students/1/", "/student/1"} {
		if _, matched := MatchRoutePattern("/students/{id}", path); matched {
			t.Errorf("%s must not Match", path)
		}
	}
}
