package roc

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"linker"
)

type Robot struct {
	gbot    *gobot.Gobot
	control *gobot.Robot
}

// TODO config file ?
type Roc struct {
	Robot
	arduino *firmata.FirmataAdaptor
	led     *gpio.LedDriver
}

func NewRoc(link *linker.Linker) *Roc {

	roc := &Roc{}
	roc.gbot = gobot.NewGobot()
	roc.apiCreate()

	//TODO config file
	roc.arduino = firmata.NewFirmataAdaptor("arduino", "/dev/ttyACM0")
	roc.led = gpio.NewLedDriver(roc.arduino, "led", "13")
	work := func() {
		for {
			select {
			case b := <-link.In:
				switch b {
				case 0x01:
					roc.forward(255)
				case 0x00:
					roc.stop()
				}
			}
		}
	}
	roc.control = roc.gbot.AddRobot(gobot.NewRobot("control",
		//		[]gobot.Connection{roc.arduino},
		//		[]gobot.Device{roc.led},
		work,
	))
	roc.controlBind()
	return roc
}

func (roc *Roc) Start() error {
	roc.gbot.Start()
	return nil
}

func (roc *Roc) Stop() error {
	return roc.gbot.Robots().Stop()[0]
}
