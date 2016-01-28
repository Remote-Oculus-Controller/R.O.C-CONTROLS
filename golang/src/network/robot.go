package network

import (
	"fmt"
	"net"
)

func RobotLink(ip, port string, in chan byte) {

	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		select {
		case b := <-in:
			//TODO buff channel
			msg := make([]byte, 2)
			msg[0] = 0x0A
			msg[1] = b
			conn.Write(msg)
		}
	}
}
