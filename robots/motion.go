package robots

import (
	"github.com/Happykat/R.O.C-CONTROLS"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

type Motion struct {
	*roc.RocRobot
	arduino        *firmata.FirmataAdaptor
	servoX, servoY *gpio.ServoDriver
	x, y           uint8
}

const (
	M_TAG = 0x10
	RTX   = 0x11
	RTY   = 0x12
)

func NewMotion() *Motion {

	m := new(Motion)
	m.RocRobot = roc.NewRocRobot(nil)
	m.arduino = firmata.NewFirmataAdaptor("arduino", "/dev/ttyACM0")
	m.servoX = gpio.NewServoDriver(m.arduino, "servoX", "5")
	m.servoY = gpio.NewServoDriver(m.arduino, "servoY", "6")
	/*
		m.AddFunc(m.rotateX, 4, nil, "")
		m.AddFunc(m.rotateY, 3, nil, "")
	*/
	m.Robot = gobot.NewRobot("motion",
		[]gobot.Connection{m.arduino},
		[]gobot.Device{m.servoX, m.servoY})
	m.x = 90
	m.y = 90
	return m
}

/*
func (m *Motion) rotateX(b []byte) error {
*/
/*a, _ := misc.DecodeUint8(b)*/ /*

	m.servoX.Move(m.x)
	m.x += 20
	return nil
}

func (m *Motion) rotateY(b []byte) error {
*/
/*	a, _ := misc.DecodeUint8(b)
 */ /*

	m.servoY.Move(m.y)
	m.y += 20
	return nil
}

func (m *Motion) moveAPI(params map[string]interface{}) interface{} {

	return 200
}
*/
