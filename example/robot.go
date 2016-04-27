package main

import (
)
import (
	"R.O.C-CONTROLS"
)

func main() {

	clientL := roc.NewLinker("", "127.0.0.1:4343", false, true)
	roc := roc.NewRoc(clientL.RegisterChannel(true))
	roc.Start()
}
