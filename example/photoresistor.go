package main

import (
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"time"
	"fmt"
)

func main() {

	gbot := gobot.NewGobot()
	firmataAdaptor := firmata.NewFirmataAdaptor("arduino", "/dev/ttyACM0")
	sensor := gpio.NewAnalogSensorDriver(firmataAdaptor, "sensor", "0")
	led := gpio.NewLedDriver(firmataAdaptor, "led", "8")
	button := gpio.NewButtonDriver(firmataAdaptor, "button", "9")
	work := func() {
		led.On()

		gobot.On(sensor.Event("data"), func(data interface{}) {
			brightness := uint8(
				gobot.ToScale(gobot.FromScale(float64(data.(int)), 0, 1024), 0, 255),
			)
			fmt.Println("sensor", data)
			fmt.Println("brightness", brightness)
		})

		gobot.On(button.Event("push"), func(data interface{}) {
			_ = time.Millisecond
			fmt.Println("push")
			led.On()
			data, _ = firmataAdaptor.AnalogRead("0")
			fmt.Println(data)
			data, _ = firmataAdaptor.AnalogRead("1")
			fmt.Println(data)
			data, _ = firmataAdaptor.AnalogRead("2")
			fmt.Println(data)
			data, _ = firmataAdaptor.AnalogRead("3")
			fmt.Println(data)
			data, _ = firmataAdaptor.AnalogRead("4")
			fmt.Println(data)
			data, _ = sensor.Read()
			fmt.Println(data)

		})
		gobot.On(button.Event("release"), func(data interface{}) {
			fmt.Println("release")
			led.Off()
		})
		/*
		for {
			if button.Active {
				led.On()
			} else {
				led.Off()
			}
		}
		*/
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{sensor, led, button},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()

}
