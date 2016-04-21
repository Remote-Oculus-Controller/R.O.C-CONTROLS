package main

import (
	"controller"
	//zz	"linker"
	"linker"
)

func main() {

	robotL := linker.NewLinker("", "127.0.0.1:4343", false, false)
	//device := controller.NewKeyboard(robotL)
	device := controller.NewDS3(robotL)
	//device := controller.NewXbox(robotL)
	//unityL := linker.NewLinker("unity", "127.0.0.1", "4343", true)
	device.Start()
}
