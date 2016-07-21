package robots

import (
	"fmt"
	"github.com/Happykat/R.O.C-CONTROLS/misc"
	"github.com/Happykat/R.O.C-CONTROLS/rocproto"
	"github.com/hybridgroup/gobot"
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

	var err error

	d := &DataLux{iter: 0, iterMax: 0, lux: -1, iterMaxLux: 0}
	for {
		select {
		case <-time.After(3 * time.Second):
			d.iterMax = d.iter
			fmt.Println("Max lux is at: ", d.getAngle(), " degrees")
			p := rocproto.Prepare(uint32(rocproto.AiInfo_DLIGH), rocproto.Packet_DATA, rocproto.Packet_CONTROL_SERVER, rocproto.Packet_VIDEO_CLIENT)
			p.Payload, err = rocproto.PackAny(&rocproto.Coord{Ori: d.getAngle()})
			if err != nil {
				return err
			}
			return ia.Send(p)
		case <-time.After(200 * time.Millisecond):
			v, err := ia.sensorLight.Read()
			if misc.CheckError(err, "reading light sensor in photoresistor.go", false) != nil {
				return err
			}
			d.iter += 1
			temp := gobot.ToScale(gobot.FromScale(float64(v), 0, 1024), 0, 255)
			fmt.Println("Sensor data", temp)
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
