package main

import (
	//"bufio"
	"fmt"
	// "time"
	// "golang.org/x/net/websocket"
	//"io"
	// "log"
	"net"
	//"os"
	"github.com/hybridgroup/gobot"
	// "github.com/hybridgroup/gobot/platforms/firmata"
	// "github.com/hybridgroup/gobot/platforms/gpio"
)

func main() {

	in := make(chan byte)
	startTCPServer(in)
	startRobot(in)
}

func startTCPServer(in chan byte) {

	fmt.Println("Starting server...")
	ln, err := net.Listen("tcp", ":4242")
	if err != nil {
		// handle error
		fmt.Print("Error : " + err.Error() + "\n")
		return
	}
	fmt.Println("Server Started")
	conn, err := ln.Accept()
	if err != nil {
		// handle error
		fmt.Print("Error : " + err.Error() + "\n")
		return
	}
	fmt.Println("Client connected")
	go handleConnection(conn, in)
}

func handleConnection(conn net.Conn, in chan byte) {
	buff := make([]byte, 8)
	// reader := bufio.NewReader(conn)
	for {
		l, err := conn.Read(buff)
		if l > 0 {
			fmt.Println("lenght : ", l, buff)
			for i := 0; i < l; i++ {
				in <- buff[i]
			}
		}
		if err != nil {
			fmt.Print("Error : " + err.Error() + "\n")
			conn.Close()
			return
		}
	}
}

func startRobot(in chan byte) {

	// var headX uint8 = 90
	// var headY uint8 = 180

	gbot := gobot.NewGobot()

	// firmataAdaptor := firmata.NewFirmataAdaptor("arduino", "COM3")
	// led := gpio.NewLedDriver(firmataAdaptor, "led", "13")
	// servoY := gpio.NewServoDriver(firmataAdaptor, "servoY", "11")
	// servoX := gpio.NewServoDriver(firmataAdaptor, "servoX", "6")

	work := func() {
		// gobot.Every(1*time.Second, func() {
		// 	led.Toggle()
		// })
		// servoX.Move(headY)
		// servoY.Move(headY)
		for {
			select {
			case b := <-in:
				switch b {
				case 0x0A:
					b := <-in
					// if b == 'z' && headY > 80 {
					// 	headY -= 5
					// 	servoY.Move(headY)
					// } else if b == 's' && headY < 180 {
					// 	headY += 5
					// 	servoY.Move(headY)
					// } else if b == 'q' && headX > 0 {
					// 	headX -= 5
					// 	servoX.Move(headX)
					// } else if b == 'd' && headX < 180 {
					// 	headX += 5
					// 	servoX.Move(headX)
					// } else
					if b == 1 {
						fmt.Println("forward")
					} else if b == 2 {
						fmt.Println("stop")
					}
					break
				}
				break
			}
		}
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{},
		[]gobot.Device{},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}

func command(c byte) {

}
