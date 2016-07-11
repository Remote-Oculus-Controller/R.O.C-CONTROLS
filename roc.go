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
	aiLock       chan bool
	cmd          chan *Packet
	data         chan *Packet
	error        chan *Packet
}

const (
	MAGIC = uint32(Packet_MAGIC_Number)

	Mask_DST   = uint32(Packet_MASK_DEST)
	Shift_SEND = uint32(Packet_SHIFT_SENT)
	Mask_SEND  = uint32(Packet_MASK_DEST) << Shift_SEND
	Shift_TYPE = uint32(Packet_SHIFT)
	Mask_Type  = uint32(Packet_MASK_TYPE) << Shift_TYPE

	CMD   = uint32(Packet_COMMAND) << Shift_TYPE
	DATA  = uint32(Packet_DATA) << Shift_TYPE
	ERROR = uint32(Packet_ERROR) << Shift_TYPE
)

func NewRoc(lS, rS string, lT, rT bool) *Roc {

	roc := &Roc{
		Gobot:  gobot.NewGobot(),
		cmap:   make(map[uint32]func(*Packet) error),
		l:      NewLinker(lS, rS, lT, rT),
		aiLock: make(chan bool),
		cmd:    make(chan *Packet),
		data:   make(chan *Packet),
		error:  make(chan *Packet),
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
			switch b.Header & Mask_Type {
			case CMD:
				r.cmd <- b
			case DATA:
				r.data <- b
			case ERROR:
				r.error <- b
			default:
				log.Println("Unknown Type", b.Header&Mask_Type)
			}
		}
	}
}

func (r *Roc) handleCmd(ch chan *Packet) {

	for {
		select {
		case <-r.aiLock:
		NoneLoop:
			for {
				select {
				case p := <-ch:
					p = p
				case <-r.aiLock:
					break NoneLoop
				}
			}
		case p := <-ch:
			f, k := r.cmap[p.ID]
			if k {
				err := f(p)
				if err != nil {
					log.Println(err.Error())
				}
			} else {
				log.Println("Unknow code", p.ID)
			}
		}
	}
}

func (r *Roc) handleData(ch chan *Packet) {
	p := <-ch
	log.Printf("Data ! ==> %+v", p)
}

func (r *Roc) handleError(ch chan *Packet) {
	p := <-ch
	log.Printf("Error ! ==> %+v", p)
}

func (r *Roc) Start() error {

	r.l.Start()
	r.NewAI()
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

func (roc *Roc) Stop() []error {
	return roc.Gobot.Stop()
}

func (roc *Roc) AddRocRobot(m *RocRobot) {
	if roc.Robot(m.Name) != nil {
		log.Println("Warning !", m.Name, "bot overwritten")
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
