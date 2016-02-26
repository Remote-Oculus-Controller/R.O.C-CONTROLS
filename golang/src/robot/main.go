package main

import (
	"fmt"
	"html"
	"net/http"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/api"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

func main() {
	gbot := gobot.NewGobot()

	a := api.NewAPI(gbot)
	a.AddHandler(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q \n", html.EscapeString(r.URL.Path))
	})

	gbot.AddCommand("custom_gobot_command",
		func(params map[string]interface{}) interface{} {
			return "This command is attached to the mcp!"
		})

	hello := gbot.AddRobot(gobot.NewRobot("hello"))

	hello.AddCommand("hi_there", func(params map[string]interface{}) interface{} {
		return fmt.Sprintf("This command is attached to the robot %v", hello.Name)
	})

	firmataAdaptor := firmata.NewFirmataAdaptor("arduino", "/dev/ttyACM0")
	led := gpio.NewLedDriver(firmataAdaptor, "led", "13")

	control := gbot.AddRobot(gobot.NewRobot("control",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{led},
	))

	control.AddCommand("forward", func(params map[string]interface{}) interface{} {
		return fmt.Sprintf("Command atach to the robot %v : %s", control.Name, "forward")
	})

	gbot.Start()

}
