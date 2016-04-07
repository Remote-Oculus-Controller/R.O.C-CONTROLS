package controller

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/joystick"
)

const (
	DS3_CF = "ds3_map.json"
)

type Dualshock3 struct {
	Controller
}

func (d *Dualshock3) Type() string {
	return "DS3"
}

func NewDS3(out chan byte) *Dualshock3 {

	d := Dualshock3{}
	d.out = out
	err := d.mapControl(DS3_CF)
	if err != nil {
		fmt.Println("Can't start controller. panicking...")
		panic(fmt.Sprintf(err.Error()))
	}
	joystickAdaptor := joystick.NewJoystickAdaptor("ps3")
	joystick := joystick.NewJoystickDriver(joystickAdaptor,
		"ps3",
		"./src/controller/json/dualshock3.json",
	)
	work := func() {
		gobot.On(joystick.Event("square_press"), func(data interface{}) {
			d.packet([]byte{d.cmap["square_p"]})
		})
		gobot.On(joystick.Event("square_release"), func(data interface{}) {
			d.packet([]byte{d.cmap["square_r"]})
		})
		gobot.On(joystick.Event("triangle_press"), func(data interface{}) {
			fmt.Println("triangle_press")
		})
		gobot.On(joystick.Event("triangle_release"), func(data interface{}) {
			fmt.Println("triangle_release")
		})
		gobot.On(joystick.Event("left_x"), func(data interface{}) {

		})
		gobot.On(joystick.Event("left_y"), func(data interface{}) {
			fmt.Println("left_y", data)
		})
		gobot.On(joystick.Event("right_x"), func(data interface{}) {
			fmt.Println("right_x", data)
		})
		gobot.On(joystick.Event("right_y"), func(data interface{}) {
			fmt.Println("right_y", data)
		})
	}

	d.robot = gobot.NewRobot("joystickBot",
		[]gobot.Connection{joystickAdaptor},
		[]gobot.Device{joystick},
		work,
	)
	return &d
}
