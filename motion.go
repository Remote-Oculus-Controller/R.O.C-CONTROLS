package roc

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"log"
)

type Motion struct {
	*RocRobot
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

	MOUV = 0xA4
	STOP = 0xAF

	DEFAULT_CAM_X = 90
	DEFAULT_CAM_Y = 135

	MAXSPEED = 90
)

func NewMotion() *Motion {

	m := new(Motion)
	m.RocRobot = NewRocRobot(nil)
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
		[]gobot.Device{m.servoX, m.servoY, m.motorR, m.motorL},
		work)
	m.AddFunc(m.moveCam, CAM, nil, "moveCam")
	m.AddFunc(m.getCamPos, GCAM, m.getCamPosApi, "getCamAngle")
	m.AddFunc(m.resetCam, RCAM, m.resetCamAPI, "resetCam")
	m.AddFunc(m.move, MOUV, nil, "mv")
	m.AddEvent("move")
	return m
}

func (m *Motion) moveCam(p *Packet) error {

	var g Gyro

	fmt.Println("Moving Camera")
	err := UnpackAny(p.Payload, &g)
	if err != nil {
		log.Println("Impossible conversion Message is not a Gyro")
		return err
	}
	x := uint8(gobot.ToScale(gobot.FromScale(g.X, -90, 90), 0, 180))
	y := uint8(gobot.ToScale(gobot.FromScale(g.Y, -35, 35), 90, 180))
	m.X = float64(x)
	m.Y = float64(y)
	fmt.Print(x, y, m.X, m.Y)
	m.servoX.Move(x)
	m.servoY.Move(y)
	return nil
	//	return m.getCamPos(p)
}

func (m *Motion) getCamPos(p *Packet) error {

	var err error

	ReverseTo(p, Packet_DATA)
	g := Gyro{m.X - DEFAULT_CAM_X, m.Y - DEFAULT_CAM_Y}
	p.Payload, err = PackAny(&g)
	if err != nil {
		return err
	}
	return m.Send(p)
}

func (m *Motion) getCamPosApi(params map[string]interface{}) interface{} {
	return m.Gyro
}

func (m *Motion) resetCam(p *Packet) error {

	m.servoY.Move(DEFAULT_CAM_Y)
	m.servoX.Move(DEFAULT_CAM_X)
	return nil
}

func (m *Motion) resetCamAPI(params map[string]interface{}) interface{} {

	m.resetCam(nil)
	return "Camera reset to original position"
}

func (m *Motion) move(p *Packet) error {

	n := &Mouv{}
	err := UnpackAny(p.Payload, n)
	if err != nil {
		log.Println("Impossible conversion Message is not a Mouv")
		return err
	}
	gobot.Publish(m.Event("move"), *n)
	fmt.Println("Spinning MOTORS !")
	/*
		y := math.Sin(n.Gspeed)
		x := math.Cos(n.Gspeed)
	*/
	//s := uint8(n.Gspeed)
	//m.motorL.Speed(byte(s))
	return nil
}

func (m *Motion) Equal(r *gobot.Robot) {

	m.arduino = r.Connection("arduino").(*firmata.FirmataAdaptor)
	m.servoY = r.Device("servoY").(*gpio.ServoDriver)
	m.servoX = r.Device("servoX").(*gpio.ServoDriver)
	m.motorL = r.Device("motorL").(*gpio.MotorDriver)
	m.motorR = r.Device("motorR").(*gpio.MotorDriver)
	fmt.Println("Robot src\n", r, "\nRobot Copy\n", m, "\n")
	m.Robot = r
}

func (m *Motion) moveForward() {

	m.motorL.Speed(MAXSPEED)
	m.motorR.Speed(MAXSPEED)
}

func (m *Motion) moveBackward() {

	m.motorL.Speed(MAXSPEED) // -
	m.motorR.Speed(MAXSPEED) // -
}

func (m *Motion) stopMoving() {

	m.motorL.Speed(0)
	m.motorR.Speed(0)
}

func (m *Motion) turnLeft() {

	m.motorR.Speed(MAXSPEED)
	m.motorL.Speed(MAXSPEED) // -
}

func (m *Motion) turnRight() {

	m.motorR.Speed(MAXSPEED) //-
	m.motorL.Speed(MAXSPEED)
}
