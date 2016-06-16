package main

import (
	"fmt"
	"github.com/Happykat/R.O.C-CONTROLS"
	"github.com/Happykat/R.O.C-CONTROLS/protoext"
	"github.com/Happykat/R.O.C-CONTROLS/robots"
	"github.com/golang/protobuf/proto"
	"net"
)

func main() {

	buff := make([]byte, 128)
	conn, err := net.Dial("tcp", "localhost:4343")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()
	p := roc.Packet{}
	p.Magic = roc.MAGIC
	p.Header = (uint32(roc.Packet_COMMAND) << uint32(roc.Packet_SHIFT)) | uint32(roc.Packet_CONTROL_SERVER) | (uint32(roc.Packet_VIDEO_CLIENT) << uint32(roc.Packet_SHIFT_SENT))
	p.ID = 178
	m, err := proto.Marshal(&p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(m)
	conn.Write(m)
	r, err := conn.Read(buff)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = proto.Unmarshal(buff[:r], &p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	c := robots.Coord{}
	err = protoext.UnpackAny(p.GetPayload(), &c)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Gps Coordinates : %+v", c)

}
