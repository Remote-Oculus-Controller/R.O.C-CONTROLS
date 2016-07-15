package main

import (
	"fmt"
	"github.com/Happykat/R.O.C-CONTROLS"
	"github.com/Happykat/R.O.C-CONTROLS/robots"
	"github.com/golang/protobuf/proto"
	"net"
)

func main() {

	conn, err := net.Dial("tcp", "192.168.0.9:8500")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()
	p := roc.Packet{}
	p.Magic = roc.MAGIC
	p.Header = uint32(roc.Packet_COMMAND)<<uint32(roc.Packet_SHIFT) | uint32(roc.Packet_VIDEO_SERVER)
	p.ID = 177
	c := robots.Coord{}
	c.Lat = 200
	c.Long = 150
	p.Payload, err = roc.PackAny(&c)
	m, err := proto.Marshal(&p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(m)
	conn.Write(m)
}
