package robots

import (
	"fmt"
	"github.com/Happykat/R.O.C-CONTROLS/rocproto"
	"github.com/hybridgroup/gobot"
	"log"
	"time"
)

type Data struct {
	startPushingTime time.Time
	timeDifference   time.Duration
}

func (ia *AI) pushButton(params map[string]interface{}) interface{} {

	gobot.Publish(ia.button.Event("push"), nil)
	return "button pushed"
}

func (ia *AI) releaseButton(params map[string]interface{}) interface{} {

	gobot.Publish(ia.button.Event("release"), nil)
	return "button released"
}

func (ia *AI) obstacle() {

	d := new(Data)
	ch := make(chan bool)

	gobot.On(ia.button.Event("push"), func(data interface{}) {
		d.startPushingTime = time.Now()
		log.Println("Le bouton poussoir est enfonce")
		ia.sendMessageAI(rocproto.AiInfo_OBSTACLE)
		select {
		case <-time.After(time.Second * 3):
			ia.toggle(false)
			log.Println("Ai control")
			ia.unlockRobot()
			log.Println("Ai eeleasing Control")
			ia.toggle(false)
			break
		case <-ch:
			fmt.Println("c'etait un obstacle passager, retour a la normale")
			break
		}
	})

	gobot.On(ia.button.Event("release"), func(data interface{}) {
		log.Println("Button poussoir est relache")
		ia.sendMessageAI(rocproto.AiInfo_N_OBSTACLE)
		ch <- true
	})

}

func (ia *AI) unlockRobot() {

	ia.m.moveBackward()
	<-time.After(time.Second * 2)
	ia.m.stopMoving()
	ia.m.turnLeft()
	<-time.After(time.Second * 2)
	ia.m.stopMoving()
}
