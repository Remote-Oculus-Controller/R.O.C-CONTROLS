package network

import (
	"net"
	"log"
	"github.com/Remote-Oculus-Controller/R.O.C-CONTROLS/misc"
	"github.com/golang/protobuf/proto"
)

type TcpClient struct {
	tcp  *net.TCPConn
	*RocNet
}

func NewTcpClient(r *RocNet) *TcpClient {
	return &TcpClient{RocNet: r}
}

//TODO timeout connection and try
func (l *TcpClient) Start() {

	if l.ip == "" {
		return
	}
	log.Println("Starting connection on ", l.ip)
	tcpAddr, err := net.ResolveTCPAddr("tcp", l.ip)
	misc.CheckError(err, "resolving address in linker.go/startConn", true)
	for {
		/*
		if m {
			log.Println("Listening on", tcpAddr.String())
			listener, err = net.ListenTCP("tcp", tcpAddr)
			misc.CheckError(err, "listening in linker.go/startConn", true)
			log.Println("Looking for a client...")
			l.tcp, err = listener.AcceptTCP()
			misc.CheckError(err, "Accepting client in linker.go/startCnn", true)
			log.Print("Connection acepted")
			listener.Close()
		*/
		/*
		} else {
		*/
			log.Println("Dialing...")
			l.tcp, err = net.DialTCP("tcp", nil, tcpAddr)
			misc.CheckError(err, "Dialing adresse in linker.go/startConn", true)
		log.Println("Connection acquired")
		/*
		}
		*/
		l.handleConn()
		log.Println("Closing connection")
		//listener.Close()
		l.tcp.Close()
		l.tcp = nil
	}
}

func (tcp *TcpClient) Stop() {

	tcp.tcp.Close()
	tcp.tcp = nil
	tcp.open = false
}

//TODO Insert buffer len and check
//Handle TCP connection
//
//The t parameters contain the section to which the handle is associated with Controls<=>Server Video / Controls<=>Client Videos
func (l *TcpClient) handleConn() {

	defer l.Stop()

	l.open = true
	quit := make(chan bool)
	go l.receive(quit)
	l.send(quit)
}

func (c *TcpClient) receive(quit chan bool) {

	buff := make([]byte, 128)

	defer func() { quit <- true }()
	for {
		select {
		case <-quit:
			return
		default:
			i, err := c.tcp.Read(buff[0:])
			if misc.CheckError(err, "Receiving data from conn", false) != nil {
				return
			}
			c.orderPacket(buff[:i])
		}
	}
}

func (l *TcpClient) send(quit chan bool) {

	for {
		select {
		case <-quit:
			return
		case m := <-l.out:
			b, err := proto.Marshal(m)
			if misc.CheckError(err, "tcpClient.go/send", false) != nil {
				continue
			}
			_, err = l.tcp.Write(b)
			if misc.CheckError(err, "tcpClient.go/send", false) != nil {
				quit <- true
				return
			}
		}
	}
}

func (tcp *TcpClient) Connected() bool{

	if (tcp.tcp != nil) {
		return true
	}
	return false
}