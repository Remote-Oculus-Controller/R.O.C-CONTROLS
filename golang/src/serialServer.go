package main

import (
	//"bufio"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
	//"os"
)

func WSServer(ws *websocket.Conn) {

	var msg = make([]byte, 512)
	var n int
	var err error

	fmt.Print("hello")
	if n, err = ws.Read(msg); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received : \"%s\"", msg[:n])
	if _, err = ws.Write([]byte("toto\n")); err != nil {
		log.Fatal(err)
	}
	fmt.Print("end")
}

// Echo the data received on the WebSocket.
func EchoServer(ws *websocket.Conn) {
	io.Copy(ws, ws)
}

// This example demonstrates a trivial echo server.
func main() {
	http.Handle("/echo", websocket.Handler(WSServer))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
