package main

import (
	"bufio"
	"controller"
	"fmt"
	"golang.org/x/net/websocket"
	"net"
	"strconv"
)

func main() {

	co := make(chan int)
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		// handle error
	}
	for {

		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}

		go handleConnection(conn, co)
		<-co
	}
	// for status, err := bufio.NewReader(conn).ReadString('\n'); status != "exit\n"; status, err = bufio.NewReader(conn).ReadString('\n') {
	// 	conn.Write()
	// }
}

func handleConnection(conn net.Conn, co chan int) {

	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Print(err)
		return
	}
	co <- 1
}
