package network

import (
	"github.com/Remote-Oculus-Controller/R.O.C-CONTROLS/misc"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
)

//TCPClient tcp client implementation for roc network interface
type TCPClient struct {
	tcp *net.TCPConn
	*RocNet
}

//NewTCPClient return a pointer to a TCPClient
func NewTCPClient(r *RocNet) *TCPClient {
	return &TCPClient{RocNet: r}
}

//Start client with provided network information
//Will throw a fatal error if cannot reach server
//TODO timeout connection and try
func (tcp *TCPClient) Start() {

	if tcp.ip == "" {
		return
	}
	log.Println("Starting connection on ", tcp.ip)
	tcpAddr, err := net.ResolveTCPAddr("tcp", tcp.ip)
	misc.CheckError(err, "resolving address in linker.go/startConn", true)
	for {

		log.Println("Dialing...")
		tcp.tcp, err = net.DialTCP("tcp", nil, tcpAddr)
		misc.CheckError(err, "Dialing adresse in linker.go/startConn", true)
		log.Println("Connection acquired")
		tcp.handleConn()
		log.Println("Closing connection")
		tcp.tcp.Close()
		tcp.tcp = nil
	}
}

//Stop client and free
func (tcp *TCPClient) Stop() {

	tcp.tcp.Close()
	tcp.tcp = nil
	tcp.open = false
}

//Handle TCP connection
//The t parameters contain the section to which the handle is associated with Controls<=>Server Video / Controls<=>Client Videos
func (tcp *TCPClient) handleConn() {

	defer tcp.Stop()

	tcp.open = true
	quit := make(chan bool)
	go tcp.receive(quit)
	tcp.send(quit)
}

func (tcp *TCPClient) receive(quit chan bool) {

	buff := make([]byte, 128)

	defer func() { quit <- true }()
	for {
		select {
		case <-quit:
			return
		default:
			i, err := tcp.tcp.Read(buff[0:])
			if misc.CheckError(err, "Receiving data from conn", false) != nil {
				return
			}
			tcp.orderPacket(buff[:i])
		}
	}
}

func (tcp *TCPClient) send(quit chan bool) {

	for {
		select {
		case <-quit:
			return
		case m := <-tcp.out:
			b, err := proto.Marshal(m)
			if misc.CheckError(err, "tcpClient.go/send", false) != nil {
				continue
			}
			_, err = tcp.tcp.Write(b)
			if misc.CheckError(err, "tcpClient.go/send", false) != nil {
				quit <- true
				return
			}
		}
	}
}

//Connected ...
func (tcp *TCPClient) Connected() bool {

	if tcp.tcp != nil {
		return true
	}
	return false
}
