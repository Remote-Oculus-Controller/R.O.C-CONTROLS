package roc_test

import (
	"log"
	"net/url"

	"github.com/googollee/go-engine.io/message"
	"github.com/gorilla/websocket"
)

type ClientWs struct {
	client
	conn websocket.Conn
	quit chan bool
	rcv  chan []byte
}

func (ws *ClientWs) Start(host string) {

	ws.quit = make(chan bool)
	ws.rcv = make(chan []byte)
	u := url.URL{Scheme: "ws", Host: host, Path: "/controls"}
	log.Printf("connecting to %s\n", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	go func() {
		defer ws.Stop()
		for {
			_, b, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			ws.rcv <- b
		}
	}()
	<-ws.quit
}

func (ws *ClientWs) Stop() {
	log.Println("Client Stopping")
	ws.quit <- true
}

func (ws *ClientWs) Send(b []byte) {
	ws.conn.WriteMessage(int(message.MessageBinary), b)
}

func (ws *ClientWs) Rcv() []byte {
	return <-ws.rcv
}
