package robots

import (
	"github.com/Happykat/R.O.C-CONTROLS"
	"github.com/Happykat/R.O.C-CONTROLS/protoext"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"log"
)

type Motion struct {
	*roc.RocRobot
	arduino        *firmata.FirmataAdaptor
	servoX, servoY *gpio.ServoDriver
	motorL, motorR *gpio.MotorDriver
	Gyro
}

const (
	M_TAG   = 0xA0
	CAM     = 0xA0
	GCAM    = 0xA9
	STOP    = 0xAF
	FORWARD = 0xA1
	BACK    = 0xA2
	TURN    = 0xA3
)

func NewMotion() *Motion {

	m := new(Motion)
	m.RocRobot = roc.NewRocRobot(nil)
	m.arduino = firmata.NewFirmataAdaptor("arduino", "/dev/ttyACM0")
	m.servoX = gpio.NewServoDriver(m.arduino, "servoX", "5")
	m.servoY = gpio.NewServoDriver(m.arduino, "servoY", "6")
	m.motorL = gpio.NewMotorDriver(m.arduino, "motorL", "9")
	m.motorR = gpio.NewMotorDriver(m.arduino, "motorR", "10")
	m.AddFunc(m.moveCam, CAM, nil, "")
	m.AddFunc(m.getCamPos, GCAM, m.getCamPosApi, "getCamAngle")
	work := func() {
		m.servoX.Move(uint8(m.X))
		m.servoX.Move(uint8(m.Y))
	}
	m.Robot = gobot.NewRobot("motion",
		[]gobot.Connection{m.arduino},
		[]gobot.Device{m.servoX, m.servoY},
		work)
	m.X = 90
	m.Y = 90
	return m
}

func (m *Motion) moveCam(p *roc.Packet) error {

	var g Gyro

	err := protoext.UnpackAny(p.Payload, &g)
	if err != nil {
		log.Println("Impossible conversion Message is not a Gyro")
		return err
	}
	m.servoX.Move(uint8(g.X))
	m.servoY.Move(uint8(g.Y))
	return nil
}

func (m *Motion) getCamPos(p *roc.Packet) error {

	var err error

	s := uint32(p.Header) & (uint32(roc.Packet_MASK_DEST) << uint32(roc.Packet_SHIFT_SENT))
	p.Header = (uint32(roc.Packet_DATA) << uint32(roc.Packet_SHIFT)) | s>>uint32(roc.Packet_SHIFT_SENT)
	p.Payload, err = protoext.PackAny(&m.Gyro)
	if err != nil {
		return err
	}
	return m.Send(p)
}

func (m *Motion) getCamPosApi(params map[string]interface{}) interface{} {
	return m.Gyro
}

func (m *Motion) move(p *roc.Packet) error {

	return nil
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
