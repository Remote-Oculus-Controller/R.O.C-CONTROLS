package controller

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/joystick"
	"R.O.C-CONTROLS"
)

const (
	DS3_CF = "./config/ds3_map.json"
)

type Dualshock3 struct {
	Controller
}

func (d *Dualshock3) Type() string {
	return "DS3"
}

func NewDS3(link *roc.Linker) *Dualshock3 {

	d := new(Dualshock3)
	d.link = link
	err := d.mapControl(DS3_CF)
	if err != nil {
		fmt.Println("Can't start controller. panicking...")
		panic(fmt.Sprintf(err.Error()))
	}
	joystickAdaptor := joystick.NewJoystickAdaptor("ps3")
	joystick := joystick.NewJoystickDriver(joystickAdaptor,
		"ps3",
		"./controller/json/dualshock3.json",
	)
	work := func() {
		gobot.On(joystick.Event("square_press"), func(data interface{}) {
			p := d.cmap["square_p"]
			d.packet(p.Code, p.Default)
		})
		gobot.On(joystick.Event("square_release"), func(data interface{}) {
			p := d.cmap["square_r"]
			d.packet(p.Code, p.Default)
		})
		gobot.On(joystick.Event("triangle_press"), func(data interface{}) {
			d.packet(d.cmap["triangle_p"].Code, data)
		})
		gobot.On(joystick.Event("triangle_release"), func(data interface{}) {
			d.packet(d.cmap["triangle_r"].Code, data)
		})
		gobot.On(joystick.Event("left_x"), func(data interface{}) {
			d.packet(d.cmap["left_x"].Code, data)
		})
		gobot.On(joystick.Event("left_y"), func(data interface{}) {
			d.packet(d.cmap["left_y"].Code, data)
		})
		gobot.On(joystick.Event("right_x"), func(data interface{}) {
			d.packet(d.cmap["right_x"].Code, data)
		})
		gobot.On(joystick.Event("right_y"), func(data interface{}) {
			d.packet(d.cmap["right_y"].Code, data)
		})
	}

	d.robot = gobot.NewRobot("joystickBot",
		[]gobot.Connection{joystickAdaptor},
		[]gobot.Device{joystick},
		work,
	)
	return d
}
