package roc

import (
	"errors"
	"fmt"
	"go/types"
	"log"

	"github.com/Happykat/R.O.C-CONTROLS/rocproto"
	"github.com/hybridgroup/gobot"
)

type RocRobot struct {
	*gobot.Robot
	l    *Linker
	cmap map[uint32]func(*rocproto.Packet) error
}

const PARAM_ERR = "MISSING %s in parameters"

func NewRocRobot(l *Linker) *RocRobot {

	r := new(RocRobot)
	r.l = l
	r.cmap = make(map[uint32]func(*rocproto.Packet) error)
	return r
}

func (r *RocRobot) Send(p *rocproto.Packet) error {

	p.Header = p.Header | (uint32(rocproto.Packet_CONTROL_SERVER) << uint32(rocproto.Packet_SHIFT_SEND))
	if r.l == nil {
		log.Println("Linker not set, cannot send rocproto.Packet")
		return nil
	}
	err := r.l.Send(p)
	if err != nil {
		return errors.New("Could not sent message. " + err.Error())
	}
	return nil
}

func (r *RocRobot) AddFunc(f func(*rocproto.Packet) error, code uint32, api func(map[string]interface{}) interface{}, name string) {
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

func (r *RocRobot) CheckAPIParams(m map[string]interface{}, t []types.BasicKind, params ...string) error {

	if len(t) != len(params) {
		return errors.New("Bad formating, expecting same number of type and parameters to create api function")
	}
	for _, v := range params {
		_, ok := m[v]
		assert := true
		if !ok || !assert {
			return errors.New(fmt.Sprintf(PARAM_ERR, v))
		}
	}
	return nil
}
