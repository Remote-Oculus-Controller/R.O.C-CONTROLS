package roc

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"log"
	"github.com/Happykat/R.O.C-CONTROLS/misc"
)

type Roc struct {
	*gobot.Gobot                        //Gobot
	cmap    map[byte]func([]byte) (error) //cmd func map
	l	*Linker
}

const (
	UPPERMASK  = 0xF0
	BOTTOMMASK = 0x0F
)

func NewRoc(lS, rS string, lT, rT bool) *Roc {

	roc := new(Roc)
	roc.Gobot = gobot.NewGobot()
	roc.cmap = make(map[byte]func([]byte) (error))
	roc.l = NewLinker(lS, rS, lT, rT)
	roc.apiCreate()
	roc.AddFunc(roc.forward, 1, true, "forward")
	go func() {
		for {
			select {
			case b := <-roc.l.remote.in:
				f, k := roc.cmap[b[0]]
				if k {
					f(b[1:])
				}
			}
		}
	}()
	return roc
}

func (roc *Roc) Start() error {

	roc.l.Start()
	errs := roc.Gobot.Start()
	if (errs != nil) {
		for _, err := range errs {
			log.Println(err)
		}
		panic(fmt.Sprintln("Panic starting ROC"))
	}
	return nil
}

func (roc *Roc) Stop() error {
	return roc.Stop()
}

func (roc *Roc) AddRobot(m *RocRobot) {
	if (roc.Robot(m.Name) != nil) {
		log.Println("Warning !" + m.Name + "bot overwritten")
	}
	m.l = roc.l
	for k, v := range m.cmap {
		_, ok := roc.cmap[k]
		if (ok) {
			log.Println("command code", k, "already exist skipping")
			continue
		}
		roc.cmap[k] = v
	}
	roc.Gobot.AddRobot(m.Robot)
}

//Directly add func with code, if specified create the api entry
func (roc *Roc) AddFunc(f func([]byte) (error), code byte, api bool, name string) {
	log.Println("Assigning function", name, "to code", code)
	_, k := roc.cmap[code]
	if !k {
		roc.cmap[code] = f
		if api {
			log.Println("Creating api entry for function")
			roc.AddCommand(name, func(params map[string]interface{}) interface{} {
				d, k := params["packet"]
				if k {
					v, err := misc.EncodeBytes(d)
					if err != nil {
						return fmt.Sprintln("API error:", err.Error())
					}
					err = f(v)
					return fmt.Sprintln(err.Error())
				}
				return fmt.Sprintln("API error: wrong parameter", k)
			})
		}
	} else {
		log.Println("Code", code, "already assigned")
	}
}
