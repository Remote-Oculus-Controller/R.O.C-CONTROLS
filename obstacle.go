package roc

import (
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"time"
	"fmt"
	"log"
)

type Data struct {
	startPushingTime 	int
	timeDifference		int
}

const (
	TIME 	= float64(0.5)
)

func (ia *AI) obstacle() {

	gbot := gobot.NewGobot()
	firmataAdaptor := firmata.NewFirmataAdaptor("arduino", "/dev/ttyACM0")
	button := gpio.NewButtonDriver(firmataAdaptor, "button", "13")


	work := func() {
		d := Data{}

		gobot.On(button.Event("push"), func(data interface{}) {
			d.startPushingTime = time.Now()
			log.Println("button1 pushed")
			ia.sendMessageAI("An obstacle prevents the robot to move forward")
			ia.toggle(true)
		})

		gobot.On(button.Event("release"), func(data interface{}) {
			d.timeDifference = time.Since(d.startPushingTime)
			log.Println("Envoie d'un message a l'utilisateur, un obstacle gene l'avancer du robot")
			if d.timeDifference < 3 {
				fmt.Println("c'etait un obstacle passager, retour a la normale")
			} else {
				ia.sendMessageAI("Warning, AI is taking control")
				ia.unlockRobot()
				ia.sendMessageAI("You are taking back control")
				ia.toggle(false)
			}
		})

	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{button},
		work,
	)
	gbot.AddRobot(robot)
	gbot.Start()

}

func (ia *AI) sendMessageAI(msg string) {

	p := Prepare(LOCK, Packet_DATA, Packet_CONTROL_SERVER, Packet_VIDEO_CLIENT)
	p.Payload= PackAny(&MAI{Lock: true, Msg: msg})
	ia.Send(p)
}

func (ia *AI) unlockRobot() {

	ia.m.moveBackward()
	<-time.After(time.Second * 2)
	ia.m.stopMoving()
	ia.m.turnLeft()
	<-time.After(time.Second * 2)
	ia.m.stopMoving()
}