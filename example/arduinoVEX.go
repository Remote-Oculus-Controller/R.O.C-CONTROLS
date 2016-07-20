package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"time"
)

type Motion struct {
	*gobot.Robot
	arduino        *firmata.FirmataAdaptor
	servoX, servoY *gpio.ServoDriver
	motorL, motorR *gpio.ServoDriver
	dir            float64
}

func main() {

	gbot := gobot.NewGobot()

	m := new(Motion)

	m.arduino = firmata.NewFirmataAdaptor("arduino", "COM3")
	m.servoX = gpio.NewServoDriver(m.arduino, "servoX", "6")
	m.servoY = gpio.NewServoDriver(m.arduino, "servoY", "5")
	m.motorL = gpio.NewServoDriver(m.arduino, "motorL", "9")
	m.motorR = gpio.NewServoDriver(m.arduino, "motorR", "10")

	work := func() {
		var i uint8 = 0
		gobot.Every(100*time.Millisecond, func() {
			fmt.Println("Turning", i)
			m.motorL.Move(i)
			i++
		})
	}
	m.Robot = gobot.NewRobot("motion",
		[]gobot.Connection{m.arduino},
		[]gobot.Device{m.servoX, m.servoY, m.motorR, m.motorL},
		work)

	gbot.AddRobot(m.Robot)

	gbot.Start()

}
