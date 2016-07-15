package roc

import (
	"fmt"
	"github.com/Happykat/R.O.C-CONTROLS/misc"
	"github.com/hybridgroup/gobot"
	"log"
	"time"
)

type Data struct {
	startPushingTime time.Time
	timeDifference   time.Duration
}

func (ia *AI) obstacle() {

	d := new(Data)
	ch := make(chan bool)

	gobot.On(ia.button.Event("push"), func(data interface{}) {
		d.startPushingTime = time.Now()
		log.Println("Envoie d'un message a l'utilisateur, un obstacle gene l'avancer du robot")
		ia.sendMessageAI("An obstacle prevents the robot from moving forward")
		for {
			select {
			case <-time.After(time.Second * 3):
				if ia.button.Active {
					log.Println("Ai control")
					ia.sendMessageAI("Warning, AI is taking control")
					ia.unlockRobot()
					ia.sendMessageAI("You are taking back control")
					ia.toggle(false)
					return
				}
			case <-ch:
				fmt.Println("c'etait un obstacle passager, retour a la normale")
			}
		}
	})

	gobot.On(ia.button.Event("release"), func(data interface{}) {
		ch <- true
	})

}

func (ia *AI) sendMessageAI(msg string) {

	var err error

	p := Prepare(LOCK, Packet_DATA, Packet_CONTROL_SERVER, Packet_VIDEO_CLIENT)
	p.Payload, err = PackAny(&MAI{Lock: true, Msg: msg})
	if misc.CheckError(err, "Sending Ai message", false) != nil {
		return
	}
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
