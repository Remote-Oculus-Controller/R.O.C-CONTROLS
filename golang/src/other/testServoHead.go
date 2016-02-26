package main

import (
	"bufio"
	"fmt"
	"os"

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

	work := func() {

		var headX uint8 = 90
		var headY uint8 = 180

		reader := bufio.NewReader(os.Stdin)
		servoY.Move(headY)
		servoX.Move(headX)
		reader.Read(1)
		servoX.Move(0)
		reader.Read(1)

	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{led, servoY, servoX},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
