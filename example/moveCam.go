package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Remote-Oculus-Controller/proto"
	"github.com/Remote-Oculus-Controller/proto/go"
	"github.com/golang/protobuf/proto"
)

func main() {

	var err error
	var buff [128]byte

	conn, err := net.Dial("tcp", ":8001")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	p := rocproto.Cam{X: 120, Y: 30}
	r := goPack.Prepare(rocproto.Cam_mv, rocproto.Packet_COMMAND, rocproto.Packet_VIDEO_CLIENT, rocproto.Packet_CONTROL_SERVER)
	r.Cam = p
	if err != nil {
		log.Fatal(err.Error())
	}
	b, err := proto.Marshal(r)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Sending")
	conn.Write(b)
	<-time.After(time.Second)
	i, err := conn.Read(buff[0:])
	if err != nil {
		fmt.Println(err)
		return
	}
	proto.Unmarshal(buff[:i], r)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("p: %+v\n", *r)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("pos %v - %v\n", r.Cam.X, r.Cam.Y)
}
