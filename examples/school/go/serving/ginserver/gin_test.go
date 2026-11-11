package ginserver

import (
	"testing"

	"github.com/cangrejometralleta/dove-canon/examples/school/go/serving/adaptertest"
	"github.com/cangrejometralleta/dove-canon/examples/school/go/transport"
)

var _ transport.Server = GinServer{}

func TestGinServerAnswersTheSharedContract(t *testing.T) {
	adaptertest.CheckAdapterAnswersTheSameWay(t, GinServer{}.BuildHandler)
}
