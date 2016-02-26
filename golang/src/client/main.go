package main

import (
	//"controller"
	//"linker"
	"controller"
	"misc"
)

func main() {

	device := new(controller.Dualshock3)
	device.Init()
	run := misc.Run{}
	run.Runners = append(run.Runners, device)
	misc.PrintRunners("runners", run.Runners)
	run.Start()
}
