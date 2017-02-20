package network

import (
	"errors"
	"fmt"
	"github.com/Remote-Oculus-Controller/R.O.C-CONTROLS/pkt"
	"github.com/Remote-Oculus-Controller/proto"
	"github.com/Remote-Oculus-Controller/proto/go"
	"github.com/golang/protobuf/proto"
	"log"
)

//RocNetI define method to for network implementation
type RocNetI interface {
	Start()
	Stop()
	Send(*rocproto.Packet)
	Type() rocproto.Packet_Section
	Connected() bool
}

//RocNet basic type for network implementation, have a pointer to cycle trough network interfaces and route data
type RocNet struct {
	ip      string
	open    bool
	out, in chan *rocproto.Packet
	t       rocproto.Packet_Section
	tls     bool
	next    *RocNet
}

//LRocNet typdef on RocNetI array
type LRocNet []RocNetI

//NewRocNet fill out a EocNet structure to be used by network implementation
func NewRocNet(ip string, t rocproto.Packet_Section, tls bool) *RocNet {
	return &RocNet{
		ip:   ip,
		open: false,
		t:    t,
		tls:  tls,
		in:   make(chan *rocproto.Packet, 100),
		out:  make(chan *rocproto.Packet, 512)}
}

//Append network interface to node
func (c *RocNet) Append(d *RocNet) {
	c.next = d
}

//SetInChan set the channel on which data will be pushed if valid
func (c *RocNet) SetInChan(in chan *rocproto.Packet) {
	c.in = in
}

//Connected ...
func (c *RocNet) Connected() bool {
	return c.open
}

//Type getter
func (c *RocNet) Type() rocproto.Packet_Section {
	return c.t
}

//Send push packet on outer channel to be sent
func (c *RocNet) Send(m *rocproto.Packet) {
	c.out <- m
}

func (c *RocNet) orderPacket(buff []byte) {

	m, err := checkBuffer(buff)
	if err != nil {
		c.out <- pkt.Error(rocproto.Error_Packet, err.Error(), int32(c.t&^rocproto.Packet_CONTROL_SERVER))
		return
	}
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

func (c *RocNet) route(m *rocproto.Packet) {
	if m.Header&uint32(c.t) != 0 {
		log.Println("Accepted", m)
		c.in <- m
	}
	go func() {
		for n := c.next; n != c; n = n.next {
			check := (m.Header&uint32(rocproto.Packet_MASK_DEST))&^uint32(c.t) != 0
			co := n.Connected()
			if check && co == true {
				log.Println("Packet routed")
				n.out <- m
			}
		}
	}()
}

//Send collection method, send the packet to every interface with correct destination
func (l *LRocNet) Send(m *rocproto.Packet) {
	for _, v := range *l {
		if m.Header&uint32(v.Type()) != 0 && v.Connected() == true {
			v.Send(m)
		}
	}
}

//Start all the interface in the collection
func (l *LRocNet) Start() {

	for _, v := range *l {
		go v.Start()
	}
}