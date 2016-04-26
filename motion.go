package roc

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"time"
	"github.com/hybridgroup/gobot/platforms/firmata"
)

type Motion struct {
	gobot.Robot

	mLCam	gpio.ServoDriver
	mRCam	gpio.ServoDriver
	tempSens gpio.AnalogSensorDriver
	potSens gpio.AnalogSensorDriver
	mLWheel gpio.MotorDriver
	mRWheel gpio.MotorDriver

	piezo 	gpio.LedDriver
	button1 gpio.LedDriver
	button2 gpio.LedDriver
	button3 gpio.LedDriver
	tilt	gpio.LedDriver

	moving 	bool
}

const (
	VSTOP 	= 0
	VRUN 	= 200
	ROTSERVO = 180
)

func newMotion() *gobot.Robot{

	r := Motion{}
	firmataAdaptor := firmata.NewFirmataAdaptor("arduino", "/dev/ttyACM0")

	r.mLCam = gpio.NewServoDriver(firmataAdaptor, "cameraMotor", "0")
	r.mRCam = gpio.NewServoDriver(firmataAdaptor, "cameraMotor", "1")
	r.tempSens = gpio.NewAnalogSensorDriver(firmataAdaptor, "temperature", "2")
	r.potSens = gpio.NewAnalogSensorDriver(firmataAdaptor, "potensiometer", "3")
	r.mLWheel = gpio.NewMotorDriver(firmataAdaptor, "wheelMotor", "4")
	r.mRWheel = gpio.NewMotorDriver(firmataAdaptor, "wheelMotor", "5")

	r.piezo = gpio.NewLedDriver(firmataAdaptor, "piezzo", "0")
	r.button1 = gpio.NewLedDriver(firmataAdaptor, "button", "1")
	r.button2 = gpio.NewLedDriver(firmataAdaptor, "button", "2")
	r.button3 = gpio.NewLedDriver(firmataAdaptor, "button", "3")
	r.tilt = gpio.NewLedDriver(firmataAdaptor, "tilt", "4")

	work := func() {
		gobot.On(r.motion.temperaturSensor.Event("data"), func(data interface{}) {
			temp := uint8(
				gobot.ToScale(gobot.FromScale(float64(data.(int)), 0, 1024), 0, 255),
			)
			fmt.Println("temperature", temp)

		})

		//Start and stop moving robot
		gobot.On(r.button1.Event("push"), func(data interface{}) {
			//r.moving ? moveForwardWheel(r, VRUN) : moveForwardWheel(r, VSTOP)
			if (r.moving) {
				r.mRCam.Move(ROTSERVO)
				r.mLCam.Move(0)
				//moveForwardWheel(r, VRUN)
				simulMotor(r, VRUN)
			}
			if (!r.moving) {
				//stopWheel(r)
				stopSimulMotor()
			}
		})

		//simulate pression sensor / obstacle non-mouvant
		gobot.On(r.button2.Event("push"), func(data interface{}) {
			stopSimulMotor()
			extractRobot(r)
		})
	}

	robot := gobot.NewRobot("motion",
		[]gobot.Connection{},
		[]gobot.Device{},
		work)

	return robot
}

func extractRobot(r Motion) {
	r.mLCam.Move(0)
	simulMotor(r, VRUN)
}

func simulMotor(r Motion, speed byte) {
	i := uint8(0)
	gobot.Every((speed/10)*time.Millisecond, func() {
		i += 1
		fmt.Println("Turning", i)
		r.mLCam.Move(i)
		r.mRCam.Move(ROTSERVO - i)
	})
}

func stopSimulMotor() {
	// Stop gobot.Every loop or do an another loop
}

func moveBackwardWheel(r Motion, speed byte) {

}

func moveForwardWheel(r Motion, speed byte) {
	r.mRWheel.speed(speed)
	r.mLWheel.speed(speed)
}

func stopWheel(r Motion) {
	r.mRWheel.speed(VSTOP)
	r.mLWheel.speed(VSTOP)
}