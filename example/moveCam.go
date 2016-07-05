package main

import (
	"github.com/Happykat/R.O.C-CONTROLS"
	"github.com/Happykat/R.O.C-CONTROLS/robots"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

type Motion struct {
	*roc.RocRobot
	arduino        *firmata.FirmataAdaptor
	servoX, servoY *gpio.ServoDriver
	motorL, motorR *gpio.MotorDriver
}

func main() {

	g := gobot.NewGobot()

	m := robots.NewMotion()

	g.AddRobot(m.Robot)
	g.Start()
}
