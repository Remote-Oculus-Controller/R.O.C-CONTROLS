package main

import (
	"fmt"
	"github.com/Remote-Oculus-Controller/R.O.C-CONTROLS/misc"
	"log"
	"net"
)

func main() {

	var listener *net.TCPListener
	var tcp *net.TCPConn

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":8000")
	misc.CheckError(err, "resolving address in linker.go/startConn", true)
	for {
		log.Println("Listening on", tcpAddr.String())
		listener, err = net.ListenTCP("tcp", tcpAddr)
		misc.CheckError(err, "listening in linker.go/startConn", true)
		log.Println("Looking for a client...")
		tcp, err = listener.AcceptTCP()
		misc.CheckError(err, "Accepting client in linker.go/startCnn", true)
		log.Print("Connection acepted")
		buff := make([]byte, 126)
		for {
			_, err := tcp.Read(buff[0:])
			if err != nil {
				break
			}
			fmt.Println(buff)
		}
		listener.Close()
	}
}
