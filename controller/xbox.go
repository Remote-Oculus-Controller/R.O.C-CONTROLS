package controller

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/joystick"
	"github.com/Happykat/R.O.C-CONTROLS"
)

var XBOX_CF = CF_DIR + "xbox_map.json"

type Xbox struct {
	Controller
}

func (d *Xbox) Type() string {
	return "Xbox"
}

func NewXbox(conn *roc.Linker) *Xbox {

	x := new(Xbox)
	x.link = conn
	err := x.mapControl(XBOX_CF)
	if err != nil {
		fmt.Println("Can't start controller. panicking...")
		panic(fmt.Sprintf(err.Error()))
	}
	joystickAdaptor := joystick.NewJoystickAdaptor("ps3")
	joystick := joystick.NewJoystickDriver(joystickAdaptor,
		"xbox",
		"./controller/json/Xbox.json",
	)
	work := func() {
		gobot.On(joystick.Event("a_press"), func(data interface{}) {
			p := x.cmap["a_p"]
			x.packet(p.Code, p.Default)
		})
		gobot.On(joystick.Event("a_release"), func(data interface{}) {
			p := x.cmap["a_r"]
			x.packet(p.Code, p.Default)
		})
		gobot.On(joystick.Event("left_x"), func(data interface{}) {
			x.packet(x.cmap["left_x"].Code, data)
		})
		gobot.On(joystick.Event("left_y"), func(data interface{}) {
			x.packet(x.cmap["left_y"].Code, data)
		})
		gobot.On(joystick.Event("right_x"), func(data interface{}) {
			x.packet(x.cmap["right_x"].Code, data)
		})
		gobot.On(joystick.Event("right_y"), func(data interface{}) {
			x.packet(x.cmap["right_y"].Code, data)
		})
	}

	x.robot = gobot.NewRobot("joystickBot",
		[]gobot.Connection{joystickAdaptor},
		[]gobot.Device{joystick},
		work,
	)
	return x
}
