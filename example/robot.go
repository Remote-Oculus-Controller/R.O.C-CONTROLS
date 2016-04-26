package main

import (
	"linker"
	"roc"
)

func main() {

	clientL := linker.NewLinker("", "127.0.0.1:4343", false, true)
	roc := &roc.Roc{Chr: clientL.RegisterChannel(true),
		Chl: clientL.RegisterChannel(false)}
	roc.Start()
}
