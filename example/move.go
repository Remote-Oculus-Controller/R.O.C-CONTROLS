package main

import (
	"fmt"
	"github.com/Happykat/R.O.C-CONTROLS"
	"github.com/Happykat/R.O.C-CONTROLS/robots"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
)

func main() {

	var err error
	//var buff [128]byte

	conn, err := net.Dial("tcp", "192.168.0.9:8001")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	p := robots.Mouv{90.00}
	r := roc.Prepare(robots.MOUV, roc.Packet_COMMAND, roc.Packet_VIDEO_CLIENT, roc.Packet_CONTROL_SERVER)
	r.Payload, err = roc.PackAny(&p)
	if err != nil {
		log.Fatal(err.Error())
	}
	b, err := proto.Marshal(r)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Sending")
	conn.Write(b)
}
