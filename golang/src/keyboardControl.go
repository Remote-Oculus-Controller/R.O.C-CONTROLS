package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/keyboard"
	"net"
)

func main() {

	out := make(chan byte)
	startTCPClient(out)
	startDevice(out)
}

func startTCPClient(out chan byte) {

	conn, err := net.Dial("tcp", "127.0.0.1:4242")
	if err != nil {
		return
	}
	go send(out, conn)
}

func send(out chan byte, conn net.Conn) {
	for {
		select {
		case b := <-out:
			msg := make([]byte, 1)
			msg[0] = b
			conn.Write(msg)
		}
	}
}

func startDevice(out chan byte) {

	gbot := gobot.NewGobot()
	keys := keyboard.NewKeyboardDriver("keyboard")
	work := func() {
		gobot.On(keys.Event("key"), func(data interface{}) {
			key := data.(keyboard.KeyEvent)

			switch key.Key {
			case keyboard.Z:
				out <- 0xA
				out <- 'z'
				break
			case keyboard.S:
				out <- 0xA
				out <- 's'
				break
			case keyboard.Q:
				out <- 0xA
				out <- 'q'
				break
			case keyboard.D:
				out <- 0xA
				out <- 'd'
				break
			}
		})
	}
	robot := gobot.NewRobot("keyboardbot",
		[]gobot.Connection{},
		[]gobot.Device{keys},
		work,
	)
	gbot.AddRobot(robot)
	gbot.Start()
}
