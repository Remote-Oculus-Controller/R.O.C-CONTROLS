package robots

import (
	"fmt"
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
	C_TAG = 0xA0
	CAM   = 0xA0
	RCAM  = 0xA1
	GCAM  = 0xA2

	DEFAULT_CAM_X = 90
	DEFAULT_CAM_Y = 135

	M_TAG   = 0xB0
	STOP    = M_TAG | 0xF
	FORWARD = M_TAG | 0x1
	BACK    = M_TAG | 0x2
	TURN    = M_TAG | 0x3
)

func NewMotion() *Motion {

	m := new(Motion)
	m.RocRobot = roc.NewRocRobot(nil)
	m.arduino = firmata.NewFirmataAdaptor("arduino", "/dev/ttyACM0")
	m.servoX = gpio.NewServoDriver(m.arduino, "servoX", "6")
	m.servoY = gpio.NewServoDriver(m.arduino, "servoY", "5")
	m.motorL = gpio.NewMotorDriver(m.arduino, "motorL", "9")
	m.motorR = gpio.NewMotorDriver(m.arduino, "motorR", "10")
	work := func() {
		m.resetCam(nil)
	}
	m.Robot = gobot.NewRobot("motion",
		[]gobot.Connection{m.arduino},
		[]gobot.Device{m.servoX, m.servoY},
		work)
	m.AddFunc(m.moveCam, CAM, nil, "moveCam")
	m.AddFunc(m.getCamPos, GCAM, m.getCamPosApi, "getCamAngle")
	m.AddFunc(m.resetCam, RCAM, m.resetCamAPI, "resetCam")
	return m
}

func (m *Motion) moveCam(p *roc.Packet) error {

	var g Gyro

	fmt.Println("Moving")
	err := protoext.UnpackAny(p.Payload, &g)
	if err != nil {
		log.Println("Impossible conversion Message is not a Gyro")
		return err
	}
	m.X = gobot.ToScale(gobot.FromScale(g.X, -90, 90), 0, 180)
	m.Y = gobot.ToScale(gobot.FromScale(g.Y, -35, 35), 90, 180)
	m.servoX.Move(uint8(m.X))
	m.servoY.Move(uint8(m.Y))
	return m.getCamPos(p)
}

func (m *Motion) getCamPos(p *roc.Packet) error {

	var err error

	protoext.ReverseTo(p, roc.Packet_DATA)
	g := Gyro{m.X - DEFAULT_CAM_X, m.Y - DEFAULT_CAM_Y}
	p.Payload, err = protoext.PackAny(&g)
	if err != nil {
		return err
	}
	fmt.Printf("p: %+v\n", p)
	return m.Send(p)
}

func (m *Motion) getCamPosApi(params map[string]interface{}) interface{} {
	return m.Gyro
}

func (m *Motion) resetCam(p *roc.Packet) error {

	m.servoY.Move(DEFAULT_CAM_Y)
	m.servoX.Move(DEFAULT_CAM_X)
	return nil
}

func (m *Motion) resetCamAPI(params map[string]interface{}) interface{} {

	m.resetCam(nil)
	return "Camera reset to original position"
}

func (m *Motion) move(p *roc.Packet) error {

	n := &Mouv{}
	err := protoext.UnpackAny(p.Payload, n)
	if err != nil {
		log.Println("Impossible conversion Message is not a Mouv")
		return err
	}
	m.motorL.Speed(byte(n.Left))
	m.motorR.Speed(byte(n.Right))
	return nil
}
