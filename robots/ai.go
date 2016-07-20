package robots

import (
	"github.com/Happykat/R.O.C-CONTROLS"
	"github.com/Happykat/R.O.C-CONTROLS/rocproto"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

const (
	LOCK = rocproto.MAI_TAG
)

type AI struct {
	*roc.RocRobot
	m      *Motion
	button *gpio.ButtonDriver
	lock   chan bool
}

func NewAI(r *roc.Roc) *AI {

	ai := &AI{RocRobot: roc.NewRocRobot(nil)}
	ai.lock = r.AiLock
	work := func() {
		/*for {
			<-time.After(time.Second * 2)
			ai.toggle(true)
			<-time.After(time.Second * 2)
			ai.toggle(false)
		}*/
	}
	gobot.On(r.Robot("motion").Event("move"), func(d interface{}) {
		r.Robot("gps").Command("sim")(map[string]interface{}{"mv": d})
	})
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
	p := rocproto.Prepare(LOCK, rocproto.Packet_DATA, rocproto.Packet_CONTROL_SERVER, rocproto.Packet_VIDEO_CLIENT)
	p.Payload, err = rocproto.PackAny(&rocproto.MAI{Lock: true})
	if err != nil {
		return err
	}
	return ai.Send(p)
}
