package network

import (
	"fmt"
	"net"
)

func UnityLink(ip, port string, in chan byte) {

	fmt.Println("Starting unity connection...")
	ln, err := net.Listen("tcp", ip+":"+port)
	if err != nil {
		// handle error
		fmt.Print("Error : " + err.Error() + "\n")
		return
	}
	fmt.Println("Socket established")
	conn, err := ln.Accept()
	fmt.Println("Client accepted")
	if err != nil {
		// handle error
		fmt.Print("Error : " + err.Error() + "\n")
		return
	}
	go handleConnection(conn, in)
}

func handleConnection(conn net.Conn, in chan byte) {
	buff := make([]byte, 8)
	// reader := bufio.NewReader(conn)
	for {
		l, err := conn.Read(buff)
		if l > 0 {
			fmt.Println(buff)
			for i := 0; i < l; i++ {
				in <- buff[i]
			}
		}
		if err != nil {
			fmt.Print("Error : " + err.Error() + "\n")
			conn.Close()
			return
		}
	}
}
