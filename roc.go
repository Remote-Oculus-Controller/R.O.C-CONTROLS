package roc

import (
	"fmt"
	"log"

	"github.com/Remote-Oculus-Controller/proto"
	"github.com/hybridgroup/gobot"
	"github.com/Remote-Oculus-Controller/R.O.C-CONTROLS/network"
	"github.com/Remote-Oculus-Controller/R.O.C-CONTROLS/pkt"
	"github.com/Remote-Oculus-Controller/proto/go"
)

//Roc super control super of robots
type Roc struct {
	*gobot.Gobot                                         //Gobot
	cmap         map[uint32]func(*rocproto.Packet) error //cmd func map
	l	     network.LRocNet
	AiLock       chan bool
	in	     chan *rocproto.Packet
	cmd          chan *rocproto.Packet
	data         chan *rocproto.Packet
	error        chan *rocproto.Packet
}

//NewRoc create
func NewRoc(lS, rS string, lT, rT bool) *Roc {

	roc := &Roc{
		Gobot:  gobot.NewGobot(),
		cmap:   make(map[uint32]func(*rocproto.Packet) error),
		AiLock: make(chan bool),
		l:	make([]network.RocNetI, 0),
		in:    make(chan *rocproto.Packet, 1024),
		cmd:    make(chan *rocproto.Packet, 512),
		data:   make(chan *rocproto.Packet, 512),
		error:  make(chan *rocproto.Packet, 10),
	}
	ws := network.NewWsSrv(network.NewRocNet(rS, rocproto.Packet_CONTROL_SERVER|rocproto.Packet_VIDEO_CLIENT, rT))
	ws.SetInChan(roc.in)
	loc := network.NewTcpClient(network.NewRocNet(lS, rocproto.Packet_CONTROL_SERVER|rocproto.Packet_VIDEO_SERVER, lT))
	loc.SetInChan(roc.in)
	ws.Append(loc.RocNet)
	loc.Append(ws.RocNet)
	roc.l = append(roc.l, ws, loc)
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
		case b := <-r.in:
			log.Printf("Packet ==>	%v\n", b)
			switch b.Header & goPack.MASK_TYPE {
			case goPack.CMD:
				r.cmd <- b
			case goPack.DATA:
				r.data <- b
			case goPack.ERROR:
				r.error <- b
			default:
				e := pkt.Error(rocproto.Error_Packet,
					"Unknown packet Type :	"+fmt.Sprintf("%b", b.Header&goPack.MASK_TYPE),
					int32(b.Header&goPack.MASK_SEND))
				log.Println(e.Err)
				r.l.Send(e)
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
					e := pkt.Error(rocproto.Error_CMDEX, err.Error(), int32(p.Header&goPack.MASK_SEND))
					log.Println(err.Error())
					r.l.Send(e)
				}
			} else {
				e := pkt.Error(rocproto.Error_Packet, "Unknown packet CMD ID :	"+fmt.Sprint(p.ID), int32(p.Header&goPack.MASK_SEND))
				log.Println(e.Err)
				r.l.Send(e)
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

	log.Println(r.l)
	r.l.Start()
	go r.handleChannel()
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
	m.l = &r.l
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
