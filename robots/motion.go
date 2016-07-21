package robots

import (
	"fmt"
	"github.com/Happykat/R.O.C-CONTROLS"
	"github.com/Happykat/R.O.C-CONTROLS/misc"
	"github.com/Happykat/R.O.C-CONTROLS/rocproto"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"log"
	"math"
)

type Motion struct {
	*roc.RocRobot
	arduino        *firmata.FirmataAdaptor
	servoX, servoY *gpio.ServoDriver
	motorL, motorR *gpio.ServoDriver
	rocproto.Gyro
	dir float64
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

	MAXSPEED  uint8 = 180
	CALCSPEED       = 90
	STOPSPEED uint8 = 90
	BACKWARD        = 0
)

func NewMotion() *Motion {

	m := new(Motion)
	m.RocRobot = roc.NewRocRobot(nil)
	m.arduino = firmata.NewFirmataAdaptor("arduino", "/dev/ttyACM0")
	m.servoX = gpio.NewServoDriver(m.arduino, "servoX", "6")
	m.servoY = gpio.NewServoDriver(m.arduino, "servoY", "5")
	m.motorL = gpio.NewServoDriver(m.arduino, "motorL", "9")
	m.motorR = gpio.NewServoDriver(m.arduino, "motorR", "10")
	m.dir = 0
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

func (m *Motion) Stop() {
	log.Println("Motion : reset cam pos and stopping motors")
	m.resetCam(nil)
	m.motorR.Move(90)
	m.motorL.Move(90)
}

func (m *Motion) moveCam(p *rocproto.Packet) error {

	var g rocproto.Gyro

	fmt.Println("Moving Camera")
	err := rocproto.UnpackAny(p.Payload, &g)
	if err != nil {
		log.Println("Impossible conversion Message is not a Gyro")
		return err
	}
	x := uint8(gobot.ToScale(gobot.FromScale(g.X, -90, 90), 0, 180))
	y := uint8(gobot.ToScale(gobot.FromScale(g.Y, -35, 35), 90, 180))
	m.X = float64(x)
	m.Y = float64(y)
	m.servoX.Move(x)
	m.servoY.Move(y)
	return nil
}

func (m *Motion) getCamPos(p *rocproto.Packet) error {

	var err error

	rocproto.ReverseTo(p, rocproto.Packet_DATA)
	g := rocproto.Gyro{m.X - DEFAULT_CAM_X, m.Y - DEFAULT_CAM_Y}
	p.Payload, err = rocproto.PackAny(&g)
	if err != nil {
		return err
	}
	return m.Send(p)
}

func (m *Motion) getCamPosApi(params map[string]interface{}) interface{} {
	return m.Gyro
}

func (m *Motion) resetCam(p *rocproto.Packet) error {

	m.servoY.Move(DEFAULT_CAM_Y)
	m.servoX.Move(DEFAULT_CAM_X)
	return nil
}

func (m *Motion) resetCamAPI(params map[string]interface{}) interface{} {

	m.resetCam(nil)
	return "Camera reset to original position"
}

func (m *Motion) move(p *rocproto.Packet) error {

	n := &rocproto.Mouv{}
	err := rocproto.UnpackAny(p.Payload, n)
	if err != nil {
		log.Println("Impossible conversion Message is not a Mouv")
		return err
	}

	var r int64 = 50

	theta := int64(n.Angle * 180 / math.Pi)
	n.Angle = n.Angle * 180 / math.Pi
	v_a := r * (45 - theta%90) / 45        // falloff of main motor
	v_b := misc.Min(100, 2*r+v_a, 2*r-v_a) // compensation of other motor
	lR, rR := thrust(theta, v_a, v_b)
	lS := gobot.ToScale(gobot.FromScale(CALCSPEED*(float64(lR)/100), -90, 90), 0, 180)
	rS := gobot.ToScale(gobot.FromScale(CALCSPEED*(float64(rR)/100), -90, 90), 0, 180)
	n.Speed = lS + rS/2 - 90
	fmt.Printf("Mouvement ==> %+v\nSpeed ==> L : %v	R : %v\n", n, lS, lR)
	gobot.Publish(m.Event("move"), *n)
	m.motorL.Move(uint8(lS))
	m.motorR.Move(uint8(rS))
	return nil
}

func thrust(theta, v_a, v_b int64) (int64, int64) {

	if theta < -90 {
		return -v_b, -v_a
	}
	if theta < 0 {
		return -v_a, v_b
	}
	if theta < 90 {
		return v_b, v_a
	}
	return v_a, -v_b
}

func (m *Motion) Equal(r *gobot.Robot) {

	m.arduino = r.Connection("arduino").(*firmata.FirmataAdaptor)
	m.servoY = r.Device("servoY").(*gpio.ServoDriver)
	m.servoX = r.Device("servoX").(*gpio.ServoDriver)
	m.motorL = r.Device("motorL").(*gpio.ServoDriver)
	m.motorR = r.Device("motorR").(*gpio.ServoDriver)
	m.Robot = r
}

func (m *Motion) moveForward() {

	m.motorL.Move(MAXSPEED)
	m.motorR.Move(MAXSPEED)
}

func (m *Motion) moveBackward() {

	m.motorL.Move(BACKWARD) // -
	m.motorR.Move(BACKWARD) // -
}

func (m *Motion) stopMoving() {

	m.motorL.Move(STOPSPEED)
	m.motorR.Move(STOPSPEED)
}

func (m *Motion) turnLeft() {

	m.motorR.Move(BACKWARD)
	m.motorL.Move(MAXSPEED) // -
}

func (m *Motion) turnRight() {

	m.motorR.Move(MAXSPEED) //-
	m.motorL.Move(BACKWARD)
}
