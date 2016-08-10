package roc

import (
	"github.com/hybridgroup/gobot/api"
)

func (roc *Roc) apiCreate() {

	a := api.NewAPI(roc.Gobot)
	a.Start()
}
