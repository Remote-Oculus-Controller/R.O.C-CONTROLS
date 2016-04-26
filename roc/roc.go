package roc

import (
	"bytes"
	"fmt"
	"github.com/hybridgroup/gobot"
)

type Roc struct {
	gbot     *gobot.Gobot
	control  *gobot.Robot
	motion   *gobot.Robot
	cmap     map[byte]func(*bytes.Buffer) error
	Chr, Chl chan []byte
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
			case b := <-roc.Chr:
				fmt.Println(b)
				/*
					switch b {
					case 0xAF:
						//emptyChannel(roc.Ch)
						//go roc.action(buff)
					default:
						log.Println("Wrong packet")
					}
				*/
			}
		}
	}

	roc.control = gobot.NewRobot("control",
		[]gobot.Connection{},
		[]gobot.Device{},
		work)

	roc.controlBind()
	roc.gbot.AddRobot(roc.control)
	roc.gbot.AddRobot(newMotion())
	roc.gbot.Start()
	return nil
}

func (roc *Roc) Stop() error {
	return roc.gbot.Stop()[0]
}

/*func (roc *Roc) addFunc(f func(*bytes.Buffer) error, code byte, api bool, name string) {
	log.Println("Assigning function", name, "to code", code)
	if (!roc.cmap[code]) {
		roc.cmap[code] = f
	} else {
		log.Println("Code", code, "already assigned")
	}
	if api {
		log.Println("Creating api entry for function")
		roc.control.AddCommand(name, func(params map[string]interface{}) interface{} {
			d, k := params["packet"]
			if k {
				v := bytes.Buffer{misc.GetBytes(d)}
				return f(v)
			}
			return log.Println("API", name, "Wrong arguments or format")
		})
	}
}*/

/*func (roc *Roc)action(buff bytes.Buffer)  {

	if buff.Len() < 2 {
		log.Println("Error sent packet do not contains enough data")
		return
	}
	switch buff.ReadByte() {
	case 0xA:
		roc.cmap[buff.ReadByte()](buff)
	}
}*/

func emptyChannel(in chan []byte) *bytes.Buffer {

	buff := new(bytes.Buffer)
	l := <-in
	/*
		for i := byte(0); i < l; i++ {
			buff.WriteByte(<-in)
		}
	*/
	fmt.Println(l, buff.Bytes())
	return buff
}
