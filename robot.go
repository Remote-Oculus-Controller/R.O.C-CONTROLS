package roc

import (
	"github.com/hybridgroup/gobot"
	"log"
	"fmt"
	"go/types"
	"errors"
)

type RocRobot struct {
	*gobot.Robot
	l *Linker
	cmap map[byte]func([]byte) (error)
}

const 	PARAM_ERR	= "MISSING %s in parameters"

func NewRocRobot(l *Linker) *RocRobot{

	r := new(RocRobot)
	r.l = l
	r.cmap = make(map[byte]func([]byte) error)
	return r
}

func (r *RocRobot) Send(b []byte) error {
	err := r.l.Send(b)
	if err != nil {
		fmt.Print("error1")
		return errors.New("Could not sent message. "+err.Error())
	}
	return nil
}

func (r *RocRobot) AddFunc(f func([]byte) (error), code byte, api func(map[string]interface{}) interface{}, name string) {
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

func (r *RocRobot) CheckAPIParams(m map[string] interface{}, t []types.BasicKind, params ...string) error {

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
