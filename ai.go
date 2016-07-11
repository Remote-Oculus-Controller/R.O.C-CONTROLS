package roc

import (
	"github.com/Happykat/R.O.C-CONTROLS/protoext"
	"github.com/hybridgroup/gobot"
	"time"
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
		for {
			<-time.After(time.Second * 2)
			ai.toggle(true)
			<-time.After(time.Second * 2)
			ai.toggle(false)
		}
	}
	ai.Robot = gobot.NewRobot("ai", work)
	r.AddRobot(ai.Robot)
	return ai
}

func (ai *AI) toggle(b bool) error {

	var err error

	ai.lock <- b
	p := protoext.Prepare(LOCK, DATA, Packet_CONTROL_SERVER, Packet_VIDEO_CLIENT)
	p.Payload, err = protoext.PackAny(MAI{Lock: true})
	if err != nil {
		return err
	}
	return ai.Send(p)
}
