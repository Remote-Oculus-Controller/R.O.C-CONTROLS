package gpsd

import (
	"encoding/json"
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/larsth/go-gpsdjson"
	"log"
	"time"
)

type GpsdDriver struct {
	name     string
	halt     chan bool
	pause    chan bool
	interval time.Duration
	r        GpsdReader
	w        GpsdWriter
	gobot.Eventer
	gobot.Commander
}

func NewGpsdDriver(adaptor *GpsdAdaptor, name string, t ...time.Duration) *GpsdDriver {

	gpsd := &GpsdDriver{
		name:      name,
		halt:      make(chan bool),
		pause:     make(chan bool),
		interval:  time.Second,
		r:         adaptor,
		w:         adaptor,
		Eventer:   gobot.NewEventer(),
		Commander: gobot.NewCommander(),
	}
	if len(t) > 0 {
		gpsd.interval = t[0]
	}
	gpsd.AddEvent(TPV)
	gpsd.AddEvent(ERROR)
	return gpsd
}

func (gpsd *GpsdDriver) Start() (errs []error) {

	gpsd.w.GpsdWrite(START)
	go func() {
		var tpv gpsdjson.TPV

		for {
			if line, err := gpsd.r.GpsdRead(); err == nil {
				json.Unmarshal([]byte(line), &tpv)
				if tpv.Class == TPV {
					gobot.Publish(gpsd.Event(TPV), tpv)
				}
			} else {
				log.Println("Error reading on gpsd socket", err.Error())
				gobot.Publish(gpsd.Event(ERROR), err)
				return
			}
			select {
			case <-time.After(gpsd.interval):
			case <-gpsd.halt:
				fmt.Printf("Halting\n")
				return
			}
		}
	}()
	return nil
}

func (gpsd *GpsdDriver) Name() string {
	return gpsd.name
}

// Connection returns the AnalogSensorDrivers Connection
func (gpsd *GpsdDriver) Connection() gobot.Connection {
	return gpsd.r.(gobot.Connection)
}

// Halt stops polling the analog sensor for new information
func (gpsd *GpsdDriver) Halt() (errs []error) {
	gpsd.halt <- true
	fmt.Printf("Halting gpsddriver\n")
	return nil
}

func (gpsd *GpsdDriver) TooglePause() {
	gpsd.pause <- true
	return
}
