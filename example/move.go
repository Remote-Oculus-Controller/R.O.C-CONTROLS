package main

import (
	"fmt"
	"github.com/Happykat/R.O.C-CONTROLS"
	"github.com/Happykat/R.O.C-CONTROLS/protoext"
	"github.com/Happykat/R.O.C-CONTROLS/robots"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
	"time"
)

func main() {

	var err error

	conn, err := net.Dial("tcp", "192.168.0.9:8001")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	p := robots.Gyro{X: 0, Y: 0}
	r := roc.Packet{}
	r.Magic = roc.MAGIC
	r.ID = robots.CAM
	r.Header = uint32(roc.Packet_CONTROL_SERVER)
	r.Payload, err = protoext.PackAny(&p)
	if err != nil {
		log.Fatal(err.Error())
	}
	b, err := proto.Marshal(&r)
	if err != nil {
		fmt.Println(err.Error())
	}
	conn.Write(b)
	<-time.After(time.Second)
	p.X = 35
	p.Y = 90
	r.Payload, err = protoext.PackAny(&p)
	if err != nil {
		log.Fatal(err.Error())
	}
	b, err = proto.Marshal(&r)
	if err != nil {
		fmt.Println(err.Error())
	}
	conn.Write(b)
}
