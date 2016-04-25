package roc

import (
	"fmt"
)

func setupMotion(roc *Roc) {

	gobot.On(roc.motion.temperaturSensor.Event("data"), func(data interface{}) {
		temp := uint8(
			gobot.ToScale(gobot.FromScale(float64(data.(int)), 0, 1024), 0, 255),
		)
		fmt.Println("temperatur", temp)

	})

	gobot.On(roc.motion.button1.Event("push"), func(data interface{}) {
		speed := byte(200)
		roc.motion.RightWheelMotor.speed(speed)
		})
	gobot.On(roc.motion.button1.Event("release"), func(data interface{}) {
		speed := byte(0)
		roc.motion.LeftMotor.speed(speed)
		})

}