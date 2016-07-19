package main

import (
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"time"
	"fmt"
)

type Data struct {
	lux 		float64
	iter		float64
	iterMaxLux 	float64
	iterMax		float64
}

const (
	COEFDIFF 	= float64(0.5)
)

func main() {

	gbot := gobot.NewGobot()
	firmataAdaptor := firmata.NewFirmataAdaptor("arduino", "/dev/ttyACM0")
	sensor := gpio.NewAnalogSensorDriver(firmataAdaptor, "sensor", "0", 200*time.Millisecond)
	led := gpio.NewLedDriver(firmataAdaptor, "led", "8")
	button := gpio.NewButtonDriver(firmataAdaptor, "button", "9")

	work := func() {
		d := &Data{}
		gobot.On(sensor.Event("data"), func(data interface{}) {
			d.iter += 1
			temp := gobot.ToScale(gobot.FromScale(float64(data.(int)), 0, 1024), 0, 255)
			fmt.Println("Sensor data", temp)
			if (diffIsCorrect(d.lux, temp) && temp > d.lux) {
					d.lux = temp
				d.iterMaxLux = d.iter
			}
		})

		gobot.On(button.Event("push"), func(data interface{}) {
			fmt.Println("button1 pushed")
			fmt.Println("check for maxi lux by turning around")
			d.iter = 0
			d.iterMax = 0
			d.lux = -1
			d.iterMaxLux = 0
		})

		gobot.On(button.Event("release"), func(data interface{}) {
			d.iterMax = d.iter
			fmt.Println("button1 released")
			fmt.Println("Step completed")
			fmt.Println("Max lux is at: ", d.getAngle(), " degrees")
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

func (d *Data) getAngle() float64 {
	angle := (360 / d.iterMax) * d.iterMaxLux
	return angle
}

func diffIsCorrect(old float64, new float64) bool {
	if (old == -1) {
		return true
	}
	if ((old - new) / old < COEFDIFF && (old - new) / old > -COEFDIFF) {
		return true
	}
	return false
}