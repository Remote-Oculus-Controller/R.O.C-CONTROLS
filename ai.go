package roc

import (
	"github.com/hybridgroup/gobot"
)

const (
	LOCK = uint32(MAI_TAG)
)

type AI struct {
	*RocRobot
	lock chan bool
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
	gobot.On(r.Robot("motion").Event("move"), func(d interface{}) {
		r.Robot("gps").Command("sim")(nil)
	})
	ai.Robot = gobot.NewRobot("ai", work)
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
