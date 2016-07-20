package main

import (
	"fmt"
	"github.com/Happykat/R.O.C-CONTROLS"
	"github.com/golang/protobuf/proto"
	"log"
	"math"
	"net"
	"time"
)

func main() {

	var err error
	//var buff [128]byte

	conn, err := net.Dial("tcp", "127.0.0.1:8001")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	for i := 1.0; i < 10; i++ {

		p := roc.Mouv{Speed: 0, Angle: math.Pi / i}
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
