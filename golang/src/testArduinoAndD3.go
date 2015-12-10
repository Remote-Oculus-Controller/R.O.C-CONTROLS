package main

import (
	"fmt"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/joystick"
)

func main() {
	gbot := gobot.NewGobot()

	firmataAdaptor := firmata.NewFirmataAdaptor("arduino", "/dev/ttyACM0")
	led := gpio.NewLedDriver(firmataAdaptor, "led", "13")
	servoY := gpio.NewServoDriver(firmataAdaptor, "servoY", "11")
	servoX := gpio.NewServoDriver(firmataAdaptor, "servoX", "6")
	motorPin1 := gpio.NewLedDriver(firmataAdaptor, "motor1", "9")
	motorPin2 := gpio.NewLedDriver(firmataAdaptor, "motor2", "2")
	motorPin3 := gpio.NewLedDriver(firmataAdaptor, "motor3", "3")
	joystickAdaptor := joystick.NewJoystickAdaptor("ps3")
	joystick := joystick.NewJoystickDriver(joystickAdaptor,
		"ps3",
		"/run/media/pheos/Media/Workspace/GitHub/R.O.C-Embeded/R.O.C-CONTROLS/golang/src/controller/json/dualshock3.json",
	)

	work := func() {

		var headX uint8 = 90
		var headY uint8 = 180

		servoY.Move(headY)
		servoX.Move(headX)
		motorPin1.On()
		gobot.On(joystick.Event("square_press"), func(data interface{}) {
			fmt.Println("square_press")
			motorPin2.On()
			motorPin3.Off()
		})
		gobot.On(joystick.Event("square_release"), func(data interface{}) {
			fmt.Println("square_release")
		})
		gobot.On(joystick.Event("triangle_press"), func(data interface{}) {
			fmt.Println("triangle_press")
			motorPin2.Off()
			motorPin3.On()
		})
		gobot.On(joystick.Event("triangle_release"), func(data interface{}) {
			fmt.Println("triangle_release")
		})
		gobot.On(joystick.Event("left_x"), func(data interface{}) {
			fmt.Println("left_x", data.(int16)/(32768/90))
			if data.(int16) < 0 && headX > 0 {
				headX -= 2
				servoX.Move(headX)
			} else if headX < 180 {
				headX += 2
				servoX.Move(headX)
			}
		})
		gobot.On(joystick.Event("left_y"), func(data interface{}) {
			fmt.Println("left_y", data.(int16)/(32768/90))
			if data.(int16) < 0 && headY > 80 {
				headY -= 2
				servoY.Move(headY)
			} else if headY < 180 {
				headY += 2
				servoY.Move(headY)
			}
		})
		gobot.On(joystick.Event("right_x"), func(data interface{}) {
			fmt.Println("right_x", data)
		})
		gobot.On(joystick.Event("right_y"), func(data interface{}) {
			fmt.Println("right_y", data)
		})
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor, joystickAdaptor},
		[]gobot.Device{led, servoY, servoX, joystick, motorPin1, motorPin2, motorPin3},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
