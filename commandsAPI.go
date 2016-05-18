package roc

import (
	"github.com/hybridgroup/gobot/api"
)

func (roc *Roc) apiCreate() {

	a := api.NewAPI(roc.Gobot)
	a.Debug()
	a.Start()
}

func (roc *Roc) controlBind() {

	//TODO Video Server call
}
