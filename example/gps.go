package main

import ()
import (
	"github.com/Happykat/R.O.C-CONTROLS/robots"
	"github.com/hybridgroup/gobot"
)

func main() {

	gobot := gobot.NewGobot()

	gpsd := robots.NewGPS()
	gobot.AddRobot(gpsd.Robot)

	gobot.Start()

}
