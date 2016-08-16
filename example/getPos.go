package main

import (
	"fmt"
	"net"

	"github.com/Remote-Oculus-Controller/proto"
	"github.com/golang/protobuf/proto"
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

		fmt.Printf("Gps Coordinates : %+v\n", p.Coord)
	}

}
