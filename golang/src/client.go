package main

import (
	//"controller"
	"controller"
	"linker"
	"misc"
)

func main() {

	robotL := linker.NewLinker("robot", "127.0.0.1", "4343", false)
	device := new(controller.Dualshock3)
	device.Init(robotL.Out)
	//unityL := linker.NewLinker("unity", "127.0.0.1", "4343", true)
	run := misc.Run{}
	run.Runners = append(run.Runners, robotL, device /*, unityL*/)
	run.Start()
}
