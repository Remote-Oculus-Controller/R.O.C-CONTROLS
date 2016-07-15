package main

import (
	"fmt"
	"github.com/Happykat/R.O.C-CONTROLS"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
	"time"
)

func main() {

	var err error
	//var buff [128]byte

	conn, err := net.Dial("tcp", ":8001")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	for i := 0.0; i < 255; i++ {

		p := roc.Mouv{i}
		r := roc.Prepare(roc.MOUV, roc.Packet_COMMAND, roc.Packet_VIDEO_CLIENT, roc.Packet_CONTROL_SERVER)
		r.Payload, err = roc.PackAny(&p)
		if err != nil {
			log.Fatal(err.Error())
		}
		b, err := proto.Marshal(r)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("Sending", i)
		conn.Write(b)
		<-time.After(time.Millisecond * 100)
	}
}
