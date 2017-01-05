package network

import (
	"github.com/Remote-Oculus-Controller/proto"
	"github.com/golang/protobuf/proto"
	"fmt"
	"log"
	"errors"
	"github.com/Remote-Oculus-Controller/R.O.C-CONTROLS/pkt"
	"github.com/Remote-Oculus-Controller/proto/go"
)

type RocNetI interface {
	Start()
	Stop()
	//Send(*rocproto.Packet)
	//Next() RocNetI
	//Connected() bool
}

type RocNet struct{
	ip      string
	open	bool
	out, in chan *rocproto.Packet
	t       rocproto.Packet_Section
	tls     bool
	next    *RocNet
}

type LRocNet []RocNetI

func NewRocNet(ip string, t rocproto.Packet_Section, tls bool) *RocNet {
	return &RocNet{
		ip: ip,
		open: false,
		t: t,
		tls: tls,
		in: make(chan *rocproto.Packet, 100),
		out: make(chan *rocproto.Packet, 512)}
}

func (c *RocNet) Start() {

}

func (c *RocNet) Stop() {

}

func (c *RocNet) Append(d *RocNet) {
	c.next = d
}

func (c *RocNet) SetInChan(in chan *rocproto.Packet) {
	c.in = in
}

func (c *RocNet) Connected() bool {
	return c.open
}

func (c *RocNet) orderPacket(buff []byte) {

	m, err := checkBuffer(buff)
	if err != nil {
		c.out <- pkt.Error(rocproto.Error_Packet, err.Error(), int32(c.t&^rocproto.Packet_CONTROL_SERVER))
		return
	}
	log.Println("Received ==>	", m)
	c.route(m)
}

func checkBuffer(buff []byte) (m *rocproto.Packet, err error) {

	log.Println(buff)
	m = &rocproto.Packet{}
	err = proto.Unmarshal(buff[0:], m)
	if err != nil {
		err = errors.New(fmt.Sprintln("Cannot Unmarshall packet : ", err.Error(),
			"\nlength ==> ", 5, "\nbuffer ==> ", buff[0:]))
		log.Println(err.Error())
		return nil, err
	}
	if m.Magic != goPack.MAGIC {
		log.Println("Wrong message format")
		return nil, errors.New("Wrong message format")
	}
	return m, nil
}

func (l *RocNet)route(m *rocproto.Packet) {
	if m.Header & uint32(l.t) != 0 {
		fmt.Println("Accepted", m)
		l.in <- m
	}
	go func() {
		for n := l.next; n != l; n = n.next {
			check := (m.Header & uint32(rocproto.Packet_MASK_DEST)) &^ uint32(l.t) != 0
			co := n.Connected()
			if check && co == true {
				log.Println("Packet routed")
				n.out <- m
			}
		}
	}()
}

func (l *LRocNet)Send(m *rocproto.Packet) {
	for _, v := range *l {
		n := v.(*RocNet)
		if m.Header & uint32(n.t) != 0 && n != nil && n.Connected() == true {
			log.Println("Packet routed")
			n.out <- m
		}
	}
}

func (l *LRocNet)Start() {

	for _, v := range *l {
		go v.Start()
	}
}
/*
func Start() {

}

f
func receive() {}*/
