package main

import (
	"log"
	"net/url"

	"fmt"

	"github.com/Remote-Oculus-Controller/proto"
	"github.com/golang/protobuf/proto"
	"github.com/googollee/go-engine.io/message"
	"github.com/gorilla/websocket"
)

func main() {

	u := url.URL{Scheme: "ws", Host: ":9000", Path: "/controls"}
	log.Printf("connecting to %s\n", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	c.WriteMessage(int(message.MessageBinary), []byte("test"))
	for {
		_, b, err := c.ReadMessage()
		if err != nil {
			fmt.Println(err)
			continue
		}
		p := rocproto.Packet{}
		err = proto.Unmarshal(b, &p)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("%+v\n", p)
	}
}
