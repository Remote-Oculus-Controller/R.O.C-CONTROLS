package main

import (
)
import (
	"R.O.C-CONTROLS"
)

func main() {

	clientL := roc.NewLinker("", "127.0.0.1:4343", false, true)
	r := roc.NewRoc(clientL.RegisterChannel(true))
	m := roc.NewMotion()
	r.SetMotion(m.G)
	//r.AddFunc(m.Forward, 2, true, "motion forward")
	r.Start()
}
