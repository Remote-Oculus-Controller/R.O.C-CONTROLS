package main

import (
	"controller"
	"network"
)

func main() {

	in := make(chan byte)
	//network.UnityLink("", "4343", in)
	go network.RobotLink("", "4242", in)
	controller.StartDS3(in)
}
