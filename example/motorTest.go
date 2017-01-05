package main

import (
	"time"
	"fmt"

	"gobot.io/x/gobot/platforms/firmata"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

func main() {

	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
	servo := gpio.NewServoDriver(firmataAdaptor, "9")

	work := func() {
		gobot.Every(1*time.Second, func() {
			i := uint8(gobot.Rand(180))
			fmt.Println("Turning", i)
			servo.Move(i)
		})
	}

	robot := gobot.NewRobot("servoBot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{servo},
		work,
	)

	robot.Start()
}
