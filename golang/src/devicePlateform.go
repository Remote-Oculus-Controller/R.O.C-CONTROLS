package main

import (
	"controller"
	"net"
	"strconv"
)

func main() {

	// conn, err := net.Dial("tcp", "localhost:8080")
	// if err != nil {
	// 	// handle error
	// 	return
	// }
	// fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		go handleConnection(conn)
	}
	// for status, err := bufio.NewReader(conn).ReadString('\n'); status != "exit\n"; status, err = bufio.NewReader(conn).ReadString('\n') {
	// 	conn.Write()
	// }
}

func handleConnection(conn net.Conn) {

	d := controller.Dualshock3{}
	end := []byte("\n")
	d.P = func(data int) {
		conn.Write(strconv.AppendInt(nil, int64(data), 10))
		conn.Write(end)
	}
	d.Start()
}
