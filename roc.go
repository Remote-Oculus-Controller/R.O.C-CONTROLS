package roc

import (
	"fmt"
	"log"

	"github.com/Remote-Oculus-Controller/proto"
	"github.com/hybridgroup/gobot"
)

//Roc super control super of robots
type Roc struct {
	*gobot.Gobot                                         //Gobot
	cmap         map[uint32]func(*rocproto.Packet) error //cmd func map
	l            *Linker
	AiLock       chan bool
	cmd          chan *rocproto.Packet
	data         chan *rocproto.Packet
	error        chan *rocproto.Packet
}

const (
	Magic = uint32(rocproto.Packet_MAGIC_Number)

	ShiftType = uint32(rocproto.Packet_SHIFT_TYPE)
	MaskType  = uint32(rocproto.Packet_MASK_TYPE) << ShiftType
	MaskSend  = uint32(rocproto.Packet_MASK_SEND)

	Cmd   = uint32(rocproto.Packet_COMMAND) << ShiftType
	Data  = uint32(rocproto.Packet_DATA) << ShiftType
	Error = uint32(rocproto.Packet_ERROR) << ShiftType
)

//NewRoc create
func NewRoc(lS, rS string, lT, rT bool) *Roc {

	roc := &Roc{
		Gobot:  gobot.NewGobot(),
		cmap:   make(map[uint32]func(*rocproto.Packet) error),
		l:      newLinker(lS, rS, lT, rT),
		AiLock: make(chan bool),
		cmd:    make(chan *rocproto.Packet, 512),
		data:   make(chan *rocproto.Packet, 512),
		error:  make(chan *rocproto.Packet, 10),
	}
	roc.apiCreate()
	return roc
}

func (r *Roc) handleChannel() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, "-> Recovered !!")
		}
	}()
	go r.handleCmd(r.cmd)
	go r.handleData(r.data)
	go r.handleError(r.error)
	for {
		select {
		case b := <-r.l.remote.in:
			log.Printf("Packet ==>	%v\n", b)
			switch b.Header & MaskType {
			case Cmd:
				r.cmd <- b
			case Data:
				r.data <- b
			case Error:
				r.error <- b
			default:
				e := NewError(rocproto.Error_Packet,
					"Unknown packet Type :	"+fmt.Sprintf("%b", b.Header&MaskType),
					int32(b.Header&MaskSend))
				log.Println(e.Err)
				r.l.send(e)
			}
		}
	}
}

func (r *Roc) handleCmd(ch chan *rocproto.Packet) {

	for {
		select {
		case <-r.AiLock:
		NoneLoop:
			for {
				select {
				case _ = <-ch:
					continue
				case <-r.AiLock:
					break NoneLoop
				}
			}
		case p := <-ch:
			f, k := r.cmap[p.ID]
			if k {
				err := f(p)
				if err != nil {
					e := NewError(rocproto.Error_CMDEX, err.Error(), int32(p.Header&MaskSend))
					log.Println(err.Error())
					r.l.send(e)
				}
			} else {
				e := NewError(rocproto.Error_Packet, "Unknown packet CMD ID :	"+fmt.Sprint(p.ID), int32(p.Header&MaskSend))
				log.Println(e.Err)
				r.l.send(e)
			}
		}
	}
}

func (r *Roc) handleData(ch chan *rocproto.Packet) {
	p := <-ch
	log.Printf("Data ! ==> %+v\n", p)
}

func (r *Roc) handleError(ch chan *rocproto.Packet) {
	p := <-ch
	log.Printf("Error ! ==> %+v\n", p)
}

//Start all component
func (r *Roc) Start() error {

	r.l.start()
	go func() {
		for {
			r.handleChannel()
		}
	}()
	errs := r.Gobot.Start()
	if errs != nil {
		for _, err := range errs {
			log.Println(err)
		}
		panic(fmt.Sprintln("Panic starting ROC"))
	}
	return nil
}

//Stop you know what
func (r *Roc) Stop() []error {
	return r.Gobot.Stop()
}

//AddRocRobot add one Robot to the collection
func (r *Roc) AddRocRobot(m *Robot) {
	if r.Robot(m.Name) != nil {
		log.Println("Warning ==>", m.Name, "bot overwritten")
	}
	m.l = r.l
	for k, v := range m.cmap {
		_, ok := r.cmap[k]
		if ok {
			log.Println("Warning ==> command code", k, "already exist skipping")
			continue
		}
		r.cmap[k] = v
	}
	r.Gobot.AddRobot(m.Robot)
}

//Directly add func with code, if specified create the api entry
func (r *Roc) AddFunc(f func(*rocproto.Packet) error, code uint32, api func(map[string]interface{}) interface{}, name string) {
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
