package main

import (
	"log"
	"net/url"

	"github.com/googollee/go-engine.io/message"
	"github.com/gorilla/websocket"
)

func main() {

	u := url.URL{Scheme: "ws", Host: "192.168.0.9:8001", Path: "/controls"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	c.WriteMessage(int(message.MessageBinary), []byte("test"))
}
