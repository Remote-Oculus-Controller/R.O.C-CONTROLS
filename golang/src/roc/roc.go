package roc

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"log"
)

type Roc struct {
	gbot    *gobot.Gobot
	control *gobot.Robot
	cmap    map[byte]func([]byte) error
	Ch      chan byte
}

const (
	UPPERMASK  = 0xF0
	BOTTOMMASK = 0x0F
)

//TODO error check
func (roc *Roc) Start() error {

	roc.gbot = gobot.NewGobot()
	roc.apiCreate()

	//TODO config file
	work := func() {
		for {
			select {
			case b := <-roc.Ch:
				switch b {
				case 0xAF:
					emptyChannel(roc.Ch)
					//roc.dispatch(link.In)
				default:
					log.Println("Wrong packet")
				}
			}
		}
	}
	roc.control = gobot.NewRobot("control",
		[]gobot.Connection{},
		[]gobot.Device{},
		work)
	roc.controlBind()
	roc.gbot.AddRobot(roc.control)
	roc.gbot.Start()
	return nil
}

func (roc *Roc) Stop() error {
	return roc.gbot.Stop()[0]
}

func (roc *Roc) addFunc(f func([]byte) error, code byte, api bool, name string) {
	roc.cmap[code] = f
	if api {
		roc.control.AddCommand(name, func(params map[string]interface{}) interface{} {
			_, k := params["packet"]
			if k {
				v := []byte(params["packet"].([]byte))
				return f(v)
			}
			return fmt.Errorf("Wrong arguments or format")
		})
	}
}

func emptyChannel(in chan byte) []byte {

	l := int(((<-in) & UPPERMASK) >> 4)
	buff := make([]byte, int(l))
	fmt.Println("len", l)
	for i := 0; i < l; i++ {
		buff[i] = <-in
	}
	fmt.Println(buff, len(buff))
	return buff
}
