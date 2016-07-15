package main

import (
	"fmt"
	"github.com/Happykat/R.O.C-CONTROLS"
	"github.com/Happykat/R.O.C-CONTROLS/robots"
	"github.com/golang/protobuf/proto"
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
	p := &roc.Packet{}
	for i := 0; i < 5; i++ {
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

		c := robots.Coord{}
		err = roc.UnpackAny(p.GetPayload(), &c)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("Gps Coordinates : %+v", c)
	}

}
