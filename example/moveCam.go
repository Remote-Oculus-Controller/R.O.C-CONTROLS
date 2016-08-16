package main

import (
	"fmt"
	"log"

	"net/url"

	"github.com/Remote-Oculus-Controller/proto"
	"github.com/Remote-Oculus-Controller/proto/go"
	"github.com/golang/protobuf/proto"
	"github.com/googollee/go-engine.io/message"
	"github.com/gorilla/websocket"
)

func main() {

	var err error

	u := url.URL{Scheme: "ws", Host: "192.168.0.9:8001", Path: "/controls"}
	log.Printf("connecting to %s\n", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	go func() {
		for i := float64(0.000); i > -1000; i-- {
			p := rocproto.Cam{X: float64(i), Y: 0}
			r := goPack.Prepare(uint32(rocproto.Cam_mv), rocproto.Packet_COMMAND, rocproto.Packet_VIDEO_CLIENT, rocproto.Packet_CONTROL_SERVER)
			r.Cam = &p
			if err != nil {
				log.Fatal(err.Error())
			}
			b, err := proto.Marshal(r)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Sending\n", r)
			conn.WriteMessage(int(message.MessageBinary), b)
		}
	}()
	for {
		r := new(rocproto.Packet)
		_, buff, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		proto.Unmarshal(buff, r)
		if err != nil {
			fmt.Println(err)
			return
		}
		if r.Cam != nil {
			fmt.Printf("pos %v - %v\n", r.Cam.X, r.Cam.Y)
		}
	}
}
