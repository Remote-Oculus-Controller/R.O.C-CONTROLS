package roc

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/firmata"
)

type Motion struct {
	G      *gobot.Robot

	mLCam	*gpio.ServoDriver
	mRCam	*gpio.ServoDriver
	tempSens *gpio.AnalogSensorDriver
	potSens *gpio.AnalogSensorDriver
	mLWheel *gpio.MotorDriver
	mRWheel *gpio.MotorDriver

	piezo 	*gpio.LedDriver
	button1 *gpio.ButtonDriver
	button2 *gpio.ButtonDriver
	button3 *gpio.ButtonDriver
	tilt	*gpio.LedDriver

	moving 	bool
}

const (
	VSTOP 	= 0
	VRUN 	= 200
	ROTSERVO = 180
)

func NewMotion() *Motion{

	r := new(Motion)
	firmataAdaptor := firmata.NewFirmataAdaptor("arduino", "/dev/ttyACM0")

	r.mLCam = gpio.NewServoDriver(firmataAdaptor, "cameraMotorL", "5")
	r.mRCam = gpio.NewServoDriver(firmataAdaptor, "cameraMotorR", "6")
	r.tempSens = gpio.NewAnalogSensorDriver(firmataAdaptor, "temperature", "2")
	r.potSens = gpio.NewAnalogSensorDriver(firmataAdaptor, "potensiometer", "3")
	/*r.mLWheel = gpio.NewMotorDriver(firmataAdaptor, "LWheel", "4")
	r.mRWheel = gpio.NewMotorDriver(firmataAdaptor, "RWheel", "5")*/

	r.piezo = gpio.NewLedDriver(firmataAdaptor, "piezzo", "0")
	r.button1 = gpio.NewButtonDriver(firmataAdaptor, "button1", "1")
	r.button2 = gpio.NewButtonDriver(firmataAdaptor, "button2", "2")
	r.button3 = gpio.NewButtonDriver(firmataAdaptor, "button3", "3")
	r.tilt = gpio.NewLedDriver(firmataAdaptor, "tilt", "4")

	work := func() {
		/*
		gobot.On(r.motion.temperaturSensor.Event("data"), func(data interface{}) {
			temp := uint8(
				gobot.ToScale(gobot.FromScale(float64(data.(int)), 0, 1024), 0, 255),
			)
			fmt.Println("temperature", temp)

		})

		*/
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


	r.G = gobot.NewRobot("motion",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{r.mLCam, r.mRCam},
		work)


	return r
}

func (m *Motion) Forward(b []byte) (byte, error)  {
	simulMotor(m, 200)
	return 200, nil
}

func extractRobot(r *Motion) {
	r.mLCam.Move(0)
	simulMotor(r, VRUN)
}

func simulMotor(r *Motion, speed byte) {
	i := uint8(gobot.Rand(180))
	fmt.Println("Turning", i)
	r.mLCam.Move(i)
	r.mRCam.Move(i)
	/*gobot.Every((time.Duration(speed)/10)*time.Millisecond, func() {
		i += 1
		fmt.Println("Turning", i)
		r.mLCam.Move(i)
		r.mRCam.Move(ROTSERVO - i)
	})*/
}

func stopSimulMotor() {
	// Stop gobot.Every loop or do an another loop
}

func moveBackwardWheel(r *Motion, speed byte) {

}

func moveForwardWheel(r *Motion, speed byte) {
	r.mRWheel.Forward(speed)
	r.mLWheel.Forward(speed)
}

func stopWheel(r Motion) {
	r.mRWheel.Halt()
	r.mLWheel.Halt()
}