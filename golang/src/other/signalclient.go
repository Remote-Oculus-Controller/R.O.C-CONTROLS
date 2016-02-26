package main

import (
	"net"
)

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:8081")
	if err != nil {
		return
	}
	var msg []byte
	msg = make([]byte, 5, 5)
	msg[0] = 1
	msg[4] = '\n'
	//str := "Message\n"
	conn.Write(msg)
	conn.Close()
}
