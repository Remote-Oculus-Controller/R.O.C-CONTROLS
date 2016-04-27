package main

import (
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"time"
	"fmt"
	"cmd/compile/internal/gc"
)

type Data struct {
	lux 		uint8
	iter		int
	iterMaxLux 	int
	iterMax		int
}

const (
	COEFDIFF 	= 0.5
)

func main() {

	gbot := gobot.NewGobot()
	firmataAdaptor := firmata.NewFirmataAdaptor("arduino", "/dev/ttyACM0")
	sensor := gpio.NewAnalogSensorDriver(firmataAdaptor, "sensor", "0", 200*time.Millisecond)
	led := gpio.NewLedDriver(firmataAdaptor, "led", "8")
	button := gpio.NewButtonDriver(firmataAdaptor, "button", "9")

	work := func() {
		Data.lux = -1
		gobot.On(sensor.Event("data"), func(data interface{}) {
			Data.iter += 1
			temp := uint8(
				gobot.ToScale(gobot.FromScale(float64(data.(int)), 0, 1024), 0, 255),
			)
			if (diffIsCorrect(temp, Data.lux)) {
				Data.lux = temp
				Data.iterMaxLux = Data.iter
			}
		})

		gobot.On(button.Event("push"), func(data interface{}) {
			fmt.Println("button1 pushed")
			fmt.Println("check for maxi lux by turning around")
			Data.iter = 0
			Data.iterMax = 0
			Data.lux = 0
			Data.iterMaxLux = 0
		})

		gobot.On(button.Event("release"), func(data interface{}) {
			Data.iterMax = Data.iter
			fmt.Println("button1 released")
			fmt.Println("Step completed")
			fmt.Println("Max lux is at: ", getAngle(), " degrees")
		})

	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{sensor, led, button},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()

}

func getAngle() int{
	angle := (360 / Data.iterMaxLux) * Data.iterMax
	return angle

}

func diffIsCorrect(old int, new int) bool {
	if (old == -1) {
		return true
	}
	if ((old - new) / old < COEFDIFF && (old - new) / old > -COEFDIFF) {
		return true
	}
	return false
}