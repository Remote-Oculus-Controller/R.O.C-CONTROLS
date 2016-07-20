package roc

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

const (
	LOCK = uint32(MAI_TAG)
)

type AI struct {
	*RocRobot
	m      *Motion
	button *gpio.ButtonDriver
	lock   chan bool
}

func (r *Roc) NewAI() *AI {

	ai := &AI{RocRobot: NewRocRobot(nil)}
	ai.lock = r.aiLock
	work := func() {
		/*for {
			<-time.After(time.Second * 2)
			ai.toggle(true)
			<-time.After(time.Second * 2)
			ai.toggle(false)
		}*/
	}
	/*	gobot.On(r.Robot("motion").Event("move"), func(d interface{}) {
		r.Robot("gps").Command("sim")(map[string]interface{}{"mv": d})
	})*/

	ai.m = NewMotion()
	ai.m.Equal(r.Robot("motion"))
	ai.button = gpio.NewButtonDriver(ai.m.arduino, "button", "13")
	ai.m.Robot.AddDevice(ai.button)
	ai.Robot = gobot.NewRobot("ai", work)
	ai.AddFunc(nil, 0, ai.pushButton, "pushButton")
	ai.AddFunc(nil, 0, ai.releaseButton, "releaseButton")
	ai.obstacle()
	r.AddRocRobot(ai.RocRobot)
	return ai
}

func (ai *AI) toggle(b bool) error {

	var err error

	ai.lock <- b
	p := Prepare(LOCK, Packet_DATA, Packet_CONTROL_SERVER, Packet_VIDEO_CLIENT)
	p.Payload, err = PackAny(&MAI{Lock: true})
	if err != nil {
		return err
	}
	return ai.Send(p)
}
