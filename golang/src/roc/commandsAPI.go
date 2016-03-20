package roc

import (
	"github.com/hybridgroup/gobot/api"
)

func (roc *Roc) apiCreate() {

	a := api.NewAPI(roc.gbot)
	a.Debug()
}

func (roc *Roc) controlBind() {

	roc.control.AddCommand("forward", func(params map[string]interface{}) interface{} {
		return roc.forward(params["speed"].(byte))
	})

	roc.control.AddCommand("stop", func(params map[string]interface{}) interface{} {
		return roc.stop()
	})
}
