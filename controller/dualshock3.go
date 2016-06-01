package controller

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/joystick"
	"os"
	"github.com/Happykat/R.O.C-CONTROLS"
)

var DS3_CF = roc.CF_DIR + "ds3_map.json"

type Dualshock3 struct {
	roc.Controller
}

func (d *Dualshock3) Type() string {
	return "DS3"
}

func NewDS3() *Dualshock3 {

	d := new(Dualshock3)
	fmt.Println(roc.CF_DIR, DS3_CF)
	err := d.MapControl(DS3_CF)
	if err != nil {
		fmt.Println("Can't start controller. panicking...")
		panic(fmt.Sprintf(err.Error()))
	}
	joystickAdaptor := joystick.NewJoystickAdaptor("ps3")
	joystick := joystick.NewJoystickDriver(joystickAdaptor,
		"ps3",
		os.Getenv("GOPATH") + "/src/github.com/Happykat/R.O.C-CONTROLS/controller/json/dualshock3.json",
	)
	work := func() {
		gobot.On(joystick.Event("square_press"), func(data interface{}) {
			fmt.Println("pressing square")
			d.Packet("s_p", data)
		})
		gobot.On(joystick.Event("square_release"), func(data interface{}) {
			d.Packet("s_r", data)
		})
		gobot.On(joystick.Event("triangle_press"), func(data interface{}) {
			d.Packet("t_p", data)
		})
		gobot.On(joystick.Event("triangle_release"), func(data interface{}) {
			d.Packet("t_r", data)
		})
		gobot.On(joystick.Event("left_x"), func(data interface{}) {
			d.Packet("left_x", data)
		})
		gobot.On(joystick.Event("left_y"), func(data interface{}) {
			d.Packet("left_y", data)
		})
		gobot.On(joystick.Event("right_x"), func(data interface{}) {
			d.Packet("right_x", data)
		})
		gobot.On(joystick.Event("right_y"), func(data interface{}) {
			d.Packet("right_y", data)
		})
	}

	d.Robot = gobot.NewRobot("joystickBot",
		[]gobot.Connection{joystickAdaptor},
		[]gobot.Device{joystick},
		work,
	)
	return d
}
