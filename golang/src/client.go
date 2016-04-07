package main

import (
	"controller"
	//zz	"linker"
	"linker"
	"misc"
)

func main() {

	robotL := linker.NewLinker("robot", "127.0.0.1", "4343", false)
	device := controller.NewKeyboard(robotL.Out)
	//device := controller.NewDS3(robotL.Out)
	//unityL := linker.NewLinker("unity", "127.0.0.1", "4343", true)
	run := misc.Run{}
	run.Runners = append(run.Runners, robotL, device /*, unityL*/)
	run.Start()
}
