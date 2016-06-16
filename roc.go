package roc

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"log"
)

type Roc struct {
	*gobot.Gobot                                //Gobot
	cmap         map[uint32]func(*Packet) error //cmd func map
	l            *Linker
}

const (
	UPPERMASK  = 0xF0
	BOTTOMMASK = 0x0F
)

func NewRoc(lS, rS string, lT, rT bool) *Roc {

	roc := new(Roc)
	roc.Gobot = gobot.NewGobot()
	roc.cmap = make(map[uint32]func(*Packet) error)
	roc.l = NewLinker(lS, rS, lT, rT)
	roc.apiCreate()
	return roc
}

func (r *Roc) handleChannel() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, "-> Recovered !!")
		}
	}()
	for {
		select {
		case b := <-r.l.remote.in:
			f, k := r.cmap[b.ID]
			if k {
				err := f(b)
				if err != nil {
					log.Println(err.Error())
				}
			} else {
				log.Println("Unknow code", b.ID)
			}
		}
	}
}

func (roc *Roc) Start() error {

	roc.l.Start()
	go func() {
		for {
			roc.handleChannel()
		}
	}()
	errs := roc.Gobot.Start()
	if errs != nil {
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
	if roc.Robot(m.Name) != nil {
		log.Println("Warning !" + m.Name + "bot overwritten")
	}
	m.l = roc.l
	for k, v := range m.cmap {
		_, ok := roc.cmap[k]
		if ok {
			log.Println("command code", k, "already exist skipping")
			continue
		}
		roc.cmap[k] = v
	}
	roc.Gobot.AddRobot(m.Robot)
}

//Directly add func with code, if specified create the api entry
func (r *Roc) AddFunc(f func(*Packet) error, code uint32, api func(map[string]interface{}) interface{}, name string) {
	if f != nil && code != 0 {
		log.Println("Assigning function", name, "to code", code)
		_, k := r.cmap[code]
		if k {
			log.Println("Code", code, "already assigned, override")
		}
		r.cmap[code] = f
	}
	if api != nil {
		log.Println("Creating api entry", name)
		r.AddCommand(name, api)
	}
}
