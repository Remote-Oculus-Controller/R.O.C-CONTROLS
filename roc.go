package roc

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"log"
	"R.O.C-CONTROLS/misc"
)

type Roc struct {
	gbot     *gobot.Gobot
	control  *gobot.Robot
	motion   *gobot.Robot
	cmap     map[byte]func([]byte) (byte, error)
	chr, chl chan []byte
}

const (
	UPPERMASK  = 0xF0
	BOTTOMMASK = 0x0F
)

func NewRoc(chr chan []byte) *Roc {

	roc := new(Roc)
	roc.chr	= chr
	roc.cmap = make(map[byte]func([]byte) (byte, error))
	work := func() {
		for {
			select {
			case b := <-roc.chr:
				f, k:= roc.cmap[b[0]]
				if (k) {
					f(b[1:])
				}
			}
		}
	}

	roc.control = gobot.NewRobot("control",
		[]gobot.Connection{},
		[]gobot.Device{},
		work)

	roc.AddFunc(roc.forward, 1, true, "forward")

	return roc
}
//TODO error check
func (roc *Roc) Start() error {

	roc.gbot = gobot.NewGobot()
//	roc.apiCreate()

	//TODO config file
	if (roc.motion != nil) {
		roc.gbot.AddRobot(roc.motion)
	}
	roc.gbot.AddRobot(roc.control)
	errs := roc.gbot.Start()
	for _, err := range errs {
		log.Println(err.Error())
	}
	return nil
}

func (roc *Roc) Stop() error {
	return roc.gbot.Stop()[0]
}

func (roc *Roc) SetMotion(m *gobot.Robot) {
	roc.motion = m
}

func (roc *Roc) AddFunc(f func([]byte) (byte, error), code byte, api bool, name string) {
	log.Println("Assigning function", name, "to code", code)
	_, k := roc.cmap[code]
	if (!k) {
		roc.cmap[code] = f
		if api {
			log.Println("Creating api entry for function")
			roc.control.AddCommand(name, func(params map[string]interface{}) interface{} {
				d, k := params["packet"]
				if k {
					v, err := misc.EncodeBytes(d)
					if err != nil {
						return fmt.Sprintln("API error:", err.Error())
					}
					r, err := f(v)
					return fmt.Sprintln(r, err.Error())
				}
				return fmt.Sprintln("API error: wrong parameter", k)
			})
		}
	} else {
		log.Println("Code", code, "already assigned")
	}
}