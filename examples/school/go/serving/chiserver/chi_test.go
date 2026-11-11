package chiserver

import (
	"testing"

	"github.com/cangrejometralleta/dove-canon/examples/school/go/serving/adaptertest"
	"github.com/cangrejometralleta/dove-canon/examples/school/go/transport"
)

var _ transport.Server = ChiServer{}

func TestChiServerAnswersTheSharedContract(t *testing.T) {
	adaptertest.CheckAdapterAnswersTheSameWay(t, ChiServer{}.BuildHandler)
}
