package main

import (
	"fmt"
	"github.com/Happykat/R.O.C-CONTROLS/rocproto"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
)

func main() {

	buff := make([]byte, 128)
	conn, err := net.Dial("tcp", "192.168.0.9:8001")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()
	p := &rocproto.Packet{}

	/*		p = rocproto.Prepare(uint32(rocproto.AiInfo_LIGHT), rocproto.Packet_COMMAND, rocproto.Packet_VIDEO_CLIENT, rocproto.Packet_CONTROL_SERVER)
			b, err := proto.Marshal(p)
			conn.Write(b)*/
	p = &rocproto.Packet{}

	for {
		r, err := conn.Read(buff)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = proto.Unmarshal(buff[:r], p)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if p.ID == 38 || p.ID == 16 {
			coor := &rocproto.Coord{}
			rocproto.UnpackAny(p.Payload, coor)
			fmt.Println(coor)
		}
		log.Printf("%+v\n", p)
	}
}
