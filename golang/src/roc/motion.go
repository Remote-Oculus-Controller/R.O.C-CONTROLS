package roc

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

type Motion struct {
	gobot.Robot

	mLCam	gpio.ServoDriver
	mRCam	gpio.ServoDriver
}

const (
	STOP = 200
)
func newMotion() *gobot.Robot{

	r := Motion{}
	r.mLCam = gpio.NewServoDriver(firmataAdaptor, "cameraMotor", "0")
	roc.motion.mRCam = gpio.NewServoDriver(firmataAdaptor, "cameraMotor", "1")
	roc.motion.temperatureSensor := gpio.NewAnalogSensorDriver(firmataAdaptor, "temperature", "2")
	roc.motion.potentiometerSensor := gpio.NewAnalogSensorDriver(firmataAdaptor, "potensiometer", "3")
	roc.motion.LeftWheelMotor := gpio.NewMotorDriver(firmataAdaptor, "wheelMotor", "4")
	roc.motion.RightWheelMotor := gpio.NewMotorDriver(firmataAdaptor, "wheelMotor", "5")

	roc.motion.piezo := gpio.NewLedDriver(firmataAdaptor, "piezzo", "0")
	roc.motion.button1 := gpio.NewLedDriver(firmataAdaptor, "button", "1")
	roc.motion.button2 := gpio.NewLedDriver(firmataAdaptor, "button", "2")
	roc.motion.button3 := gpio.NewLedDriver(firmataAdaptor, "button", "3")
	roc.motion.tilt := gpio.NewLedDriver(firmataAdaptor, "tilt", "4")

	work := func() {
		gobot.On(roc.motion.temperaturSensor.Event("data"), func(data interface{}) {
			temp := uint8(
				gobot.ToScale(gobot.FromScale(float64(data.(int)), 0, 1024), 0, 255),
			)
			fmt.Println("temperature", temp)

		})

		gobot.On(roc.motion.button1.Event("push"), func(data interface{}) {
			roc.motion.RightWheelMotor.speed(STOP)
		})
		gobot.On(roc.motion.button1.Event("release"), func(data interface{}) {
			speed := byte(0)
			roc.motion.LeftMotor.speed(speed)
		})
	}

	robot := gobot.NewRobot("motion",
		[]gobot.Connection{},
		[]gobot.Device{},
		work)

	return robot
}