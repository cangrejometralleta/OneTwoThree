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

func TestReadPatternNamesFindsEachBrace(t *testing.T) {
	names := ReadPatternNames("/a/{id}/b/{other_id}")
	if len(names) != 2 || names[0] != "id" || names[1] != "other_id" {
		t.Fatalf("names=%v", names)
	}
	if len(ReadPatternNames("/students")) != 0 {
		t.Fatal("a Pattern with no Brace Declares no Name")
	}
}
