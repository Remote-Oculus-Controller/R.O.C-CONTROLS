package main

import (
	"linker"
	"misc"
	"roc"
)

func main() {

	clientL := linker.NewLinker("robot", "127.0.0.1", "4343", true)
	roc := roc.NewRoc(clientL)
	run := misc.Run{}
	run.Runners = append(run.Runners, clientL, roc)
	misc.PrintRunners("Runners", run.Runners)
	run.Start()
}
