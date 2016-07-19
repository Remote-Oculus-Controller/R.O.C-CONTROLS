package roc

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

const (
	LOCK = uint32(MAI_TAG)
)

type AI struct {
	*RocRobot
	m      					*Motion
	buttonObstacle 	*gpio.ButtonDriver
	sensorLight 		*gpio.AnalogSensorDriver
	pending					bool
	firstTime				bool
	lock chan 			bool
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
		r.Robot("gps").Command("sim")(map[string]interface{}{"mv": d})
	})

	ai.m = NewMotion()
	fmt.Printf("Motion %+v", ai.m)
	ai.m.Equal(r.Robot("motion"))
	ai.buttonObstacle = gpio.NewButtonDriver(ai.m.arduino, "buttonObstacle", "13")
	ai.sensorLight = gpio.NewAnalogSensorDriver(ai.m.arduino, "sensorLight", "0")
	ai.m.Robot.AddDevice(ai.buttonObstacle)
	ai.m.Robot.AddDevice(ai.sensorLight)
	ai.Robot = gobot.NewRobot("ai", work)
	ai.AddFunc(nil, 0, ai.pushButton, "pushButton")
	ai.AddFunc(nil, 0, ai.releaseButton, "releaseButton")
	ai.AddFunc(nil, 0, ai.pushLightButton, "pushLightButton")
	ai.obstacle()
	ai.pending = false
	ai.firstTime = true
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
