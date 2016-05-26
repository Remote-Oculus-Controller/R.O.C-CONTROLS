package roc

import (
	"github.com/hybridgroup/gobot"
	"log"
)

type RocRobot struct {
	*gobot.Robot
	l *Linker
	cmap map[byte]func([]byte) (error)
}

func NewRocRobot(l *Linker) *RocRobot{

	r := new(RocRobot)
	r.l = l
	r.cmap = make(map[byte]func([]byte) error)
	return r
}

func (r *RocRobot) Send(b []byte) {
	r.l.Send(b)
}

func (r *RocRobot) AddFunc(f func([]byte) (error), code byte, api func(map[string]interface{}) interface{}, name string) {
	log.Println("Assigning function", name, "to code", code)
	_, k := r.cmap[code]
	if k {
		log.Println("Code", code, "already assigned, override")
	}
	r.cmap[code] = f
	if api != nil {
		log.Println("Creating api entry")
		r.AddCommand(name, api)
	}
}
