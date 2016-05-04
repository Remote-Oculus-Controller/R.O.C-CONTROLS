package main

import (
	"github.com/Happykat/R.O.C-CONTROLS/controller"
	"github.com/Happykat/R.O.C-CONTROLS"
)

//TODO opt
func main() {

	robotL := roc.NewLinker("", "127.0.0.1:4343", false, false)
	//device := controller.NewKeyboard(robotL)
	device := controller.NewDS3(robotL)
	//device := controller.NewXbox(robotL)
	//unityL := linker.NewLinker("unity", "127.0.0.1", "4343", true)
	device.Start()
}
