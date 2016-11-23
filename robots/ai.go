package robots

import (
	"github.com/Remote-Oculus-Controller/R.O.C-CONTROLS"
	"github.com/Remote-Oculus-Controller/R.O.C-CONTROLS/misc"
	"github.com/Remote-Oculus-Controller/proto"
	"github.com/Remote-Oculus-Controller/proto/go"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

type AI struct {
	*roc.Robot
	m              *Motion
	buttonObstacle *gpio.ButtonDriver
	sensorLight    *gpio.AnalogSensorDriver
	pending        bool
	firstTime      bool
	lock           chan bool
	getPos         func(map[string]interface{}) interface{}
	pattern        chan bool
}

func NewAI(r *roc.Roc) *AI {

	ai := &AI{Robot: roc.NewRocRobot(nil)}
	ai.lock = r.AiLock
	ai.pattern = make(chan bool, 2)
	work := func() {
	}
	/*gobot.On(r.Robot("motion").Event("move"), func(d interface{}) {
		r.Robot("gps").Command("sim")(map[string]interface{}{"mv": d})
	})*/
	ai.m = NewMotion()
	ai.m.Equal(r.Robot("motion"))
	//ai.getPos = r.Robot("gps").Command("getCoord")
	ai.buttonObstacle = gpio.NewButtonDriver(ai.m.arduino, "buttonObstacle", "13")
	//ai.sensorLight = gpio.NewAnalogSensorDriver(ai.m.arduino, "sensorL", "0")
	ai.m.Robot.AddDevice(ai.buttonObstacle)
	/*
		ai.m.Robot.AddDevice(ai.sensorLight)
	*/
	ai.Robot.Robot = gobot.NewRobot("ai", work)
	ai.AddFunc(nil, 0, ai.pushButton, "pushButton")
	ai.AddFunc(nil, 0, ai.stopPattern, "stopPattern")
	ai.AddFunc(nil, 0, ai.startPattern, "startPattern")
	ai.AddFunc(nil, 0, ai.releaseButton, "releaseButton")
	//ai.AddFunc(ai.startLightWorkaround, uint32(rocproto.AiCodes_LIGHT), ai.startLightDetect, "pushLightButton")
	ai.obstacle()
	ai.pending = false
	ai.firstTime = true
	return ai
}

func (ai *AI) toggle(b bool) error {

	ai.lock <- b
	if b {
		return ai.sendMessageAI(rocproto.AiCodes_LOCK)
	} else {
		return ai.sendMessageAI(rocproto.AiCodes_UNLOCK)
	}

}

func (ia *AI) sendMessageAI(id rocproto.AiCodes_Codes) error {

	var err error

	p := goPack.Prepare(uint32(id), rocproto.Packet_DATA, rocproto.Packet_CONTROL_SERVER, rocproto.Packet_VIDEO_CLIENT)
	if misc.CheckError(err, "Sending Ai message", false) != nil {
		return err
	}
	return ia.Send(p)
}
