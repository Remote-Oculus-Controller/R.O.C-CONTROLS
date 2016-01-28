package controller

import (
	"fmt"
	"github.com/hybridgroup/gobot"
)

type Controller interface {
	Type() string
	Init()
	Speed() float64
	Direction() float64
}

func Test() {
	fmt.Print("pass")
}

func StartDS3(in chan byte) {

	gbot := gobot.NewGobot()
	d3 := Dualshock3{}
	robot := d3.Create(in)
	gbot.AddRobot(robot)

	gbot.Start()
}
