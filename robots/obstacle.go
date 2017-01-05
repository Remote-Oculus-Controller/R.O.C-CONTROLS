package robots

import (
	"fmt"
	"log"
	"time"

	"github.com/Remote-Oculus-Controller/proto"
	"github.com/hybridgroup/gobot"
)

type Data struct {
	startPushingTime time.Time
	timeDifference   time.Duration
}

func (ia *AI) pushButton(params map[string]interface{}) interface{} {

	gobot.Publish(ia.Event("push"), nil)
	return "buttonObstacle pushed"
}

func (ia *AI) releaseButton(params map[string]interface{}) interface{} {

	gobot.Publish(ia.Event("release"), nil)
	return "buttonObstacle released"
}

func (ia *AI) obstacle() {

	d := new(Data)
	ch := make(chan bool)

	gobot.On(ia.Event("push"), func(data interface{}) {
		d.startPushingTime = time.Now()
		log.Println("Le bouton poussoir est enfonce")
		ia.sendMessageAI(rocproto.AiCodes_OBSTACLE)
		select {
		case <-time.After(time.Second * 3):
			ia.toggle(true)
			log.Println("Ai control")
			ia.unlockRobot()
			log.Println("Ai releasing Control")
			ia.toggle(false)
			break
		case <-ch:
			fmt.Println("c'etait un obstacle passager, retour a la normale")
			break
		}
	})

	gobot.On(ia.Event("release"), func(data interface{}) {
		log.Println("Button poussoir est relache")
		ia.sendMessageAI(rocproto.AiCodes_N_OBSTACLE)
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
