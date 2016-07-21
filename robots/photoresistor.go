package robots

import (
	"github.com/Happykat/R.O.C-CONTROLS/misc"
	"github.com/Happykat/R.O.C-CONTROLS/rocproto"
	"github.com/hybridgroup/gobot"
	"log"
	"math"
	"time"
)

type DataLux struct {
	lux        float64
	iter       float64
	iterMaxLux float64
	iterMax    float64
}

const (
	COEFDIFF = float64(0.5)
)

func (ia *AI) light() error {

	d := &DataLux{iter: 0, iterMax: 0, lux: -1, iterMaxLux: 0}
	timeout := time.After(3 * time.Second)
	tick := time.NewTicker(100 * time.Millisecond)
	ia.toggle(true)
	defer ia.toggle(false)
	for {
		select {
		case <-timeout:
			d.iterMax = d.iter
			log.Println("Max lux is at: ", d.getAngle(), " degrees")
			gobot.Publish(ia.Event("newLight"), d.getAngle())
			/*
				p.Payload, err = rocproto.PackAny(&rocproto.Coord{Ori: d.getAngle()})
				if err != nil {
					return err
				}
			*/
			tick.Stop()
			return nil
		/*
			return ia.Send(p)
		*/
		case <-tick.C:
			v, err := ia.sensorLight.Read()
			if misc.CheckError(err, "reading light sensor in photoresistor.go", false) != nil {
				return err
			}
			d.iter += 1
			temp := gobot.ToScale(gobot.FromScale(float64(v), 0, 1024), 0, 255)
			//fmt.Println("Sensor data", temp)
			if diffIsCorrect(d.lux, temp) && temp > d.lux {
				d.lux = temp
				d.iterMaxLux = d.iter
			}
		}
	}
}

func (d *DataLux) getAngle() float64 {
	angle := (360 / d.iterMax) * d.iterMaxLux
	return angle
}

func diffIsCorrect(old float64, new float64) bool {
	if old == -1 {
		return true
	}
	if (old-new)/old < COEFDIFF && (old-new)/old > -COEFDIFF {
		return true
	}
	return false
}

func (ia *AI) lightDetect(p *rocproto.Packet) error {

	var err error

	p = &rocproto.Packet{}
	p.Payload, err = rocproto.PackAny(&rocproto.Mouv{Speed: 0, Angle: math.Pi / 2})
	if misc.CheckError(err, "Packing in lightDetect", false) != nil {
		return err
	}
	ia.m.move(p)
	err = ia.light()
	if misc.CheckError(err, "Detecting light", false) != nil {
		return err
	}
	ia.m.Stop()
	return nil
}

func (ia *AI) startLightDetect(params map[string]interface{}) interface{} {

	return ia.lightDetect(nil)
}
