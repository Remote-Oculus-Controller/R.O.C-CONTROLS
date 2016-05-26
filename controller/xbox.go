package controller

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/joystick"
	"os"
	"github.com/Happykat/R.O.C-CONTROLS"
)

var XBOX_CF = roc.CF_DIR + "xbox_map.json"

type Xbox struct {
	roc.Controller
}

func (d *Xbox) Type() string {
	return "Xbox"
}

func NewXbox() *Xbox {

	x := new(Xbox)
	err := x.MapControl(XBOX_CF)
	if err != nil {
		fmt.Println("Can't start controller. panicking...")
		panic(fmt.Sprintf(err.Error()))
	}
	joystickAdaptor := joystick.NewJoystickAdaptor("ps3")
	joystick := joystick.NewJoystickDriver(joystickAdaptor,
		"xbox",
		os.Getenv("GOPATH") + "/src/github.com/Happykat/R.O.C-CONTROLS/config/xbox360_power_a_mini_proex.json",
	)
	work := func() {
		gobot.On(joystick.Event("a_press"), func(data interface{}) {
			x.Packet("a_p", data)
		})
		gobot.On(joystick.Event("a_release"), func(data interface{}) {
			x.Packet("a_r", data)
		})
		gobot.On(joystick.Event("left_x"), func(data interface{}) {
			x.Packet("left_x", data)
		})
		gobot.On(joystick.Event("left_y"), func(data interface{}) {
			x.Packet("left_y", data)
		})
		gobot.On(joystick.Event("right_x"), func(data interface{}) {
			x.Packet("right_x", data)
		})
		gobot.On(joystick.Event("right_y"), func(data interface{}) {
			x.Packet("right_y", data)
		})
	}

	x.Robot = gobot.NewRobot("joystickBot",
		[]gobot.Connection{joystickAdaptor},
		[]gobot.Device{joystick},
		work,
	)
	return x
}
