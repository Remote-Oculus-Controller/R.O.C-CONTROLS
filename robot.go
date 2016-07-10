package roc

import (
	"errors"
	"fmt"
	"github.com/hybridgroup/gobot"
	"go/types"
	"log"
)

type RocRobot struct {
	*gobot.Robot
	l    *Linker
	cmap map[uint32]func(*Packet) error
}

const PARAM_ERR = "MISSING %s in parameters"

func NewRocRobot(l *Linker) *RocRobot {

	r := new(RocRobot)
	r.l = l
	r.cmap = make(map[uint32]func(*Packet) error)
	return r
}

func (r *RocRobot) Send(p *Packet) error {

	p.Header = p.Header | (uint32(Packet_CONTROL_SERVER) << uint32(Packet_SHIFT_SENT))
	if r.l == nil {
		log.Println("Linker can not send Packet")
		return nil
	}
	fmt.Printf("Sending %+v\n", p)
	err := r.l.Send(p)
	if err != nil {
		return errors.New("Could not sent message. " + err.Error())
	}
	return nil
}

func (r *RocRobot) AddFunc(f func(*Packet) error, code uint32, api func(map[string]interface{}) interface{}, name string) {
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
		fmt.Printf("%+v", r)
		r.AddCommand(name, api)
	}
}

func (r *RocRobot) CheckAPIParams(m map[string]interface{}, t []types.BasicKind, params ...string) error {

	if len(t) != len(params) {
		return errors.New("Bad formating, expecting same number of type and parameters to create api function")
	}
	for i, v := range params {
		p, ok := m[v]
		assert := true
		fmt.Printf("%v %v %v", i, v, p)
		if !ok || !assert {
			return errors.New(fmt.Sprintf(PARAM_ERR, v))
		}
	}
	return nil
}
