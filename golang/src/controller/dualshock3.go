package controller

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/joystick"
)

const (
	ds3CF = "ds3config.json"
)

type Dualshock3 struct {
	Controller

	out chan byte
}

func (d *Dualshock3) Type() string {
	return "DS3"
}

func (d *Dualshock3) Init() {

	// if d.out == nil {
	// 	fmt.Println("Out channel for " + d.Type() + "not set")
	// 	return
	// }
	err := d.MapControl()
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
			fmt.Println("square")
			//d.out <- d.Controller.cmap["square_p"].code
		})
		gobot.On(joystick.Event("square_release"), func(data interface{}) {
			//d.out <- d.Controller.cmap["square_r"].code
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
}

func (d Dualshock3) MapControl() error {
	fmt.Println("Mapping DS3 for robot control\nStart parsing :", ds3CF)
	err := d.Controller.parseControl(ds3CF)
	if err != nil {
		fmt.Println("Failed to parse control.")
		return err
	}
	return nil
}
