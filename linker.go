package roc

import (
	"fmt"
	"net"
	"R.O.C-CONTROLS/misc"
	"log"
)

const (
	//MAGIC
	MAGIC = 0xAF

	//Type
	TYPE_SHIFT = 6
	CMD        = 1 << TYPE_SHIFT
	DATA       = (1 << 1) << TYPE_SHIFT
	ERROR      = CMD | DATA

	//Destination
	DST_SHIFT = 3
	DST_L     = 1 << DST_SHIFT
	DST_R     = (1 << 1) << DST_SHIFT
	DST_RL    = (1 << 2) << DST_SHIFT
	DST_ALL   = DST_L | DST_R | DST_RL

	//Section
	DEFAULT = 0
	MV      = 1
	SENSOR  = 1 << 1
	IA      = 1 << 2
	OTHER   = 1 << 3
	ALL     = MV | SENSOR | IA | OTHER
)

type Linker struct {
	local, remote Link
}

type Link struct {
	conn    *net.TCPConn
	out, in chan []byte
}

//TODO create start method
func NewLinker(lS, rS string, lT, rT bool) *Linker {

	var err error

	l := Linker{Link{out: make(chan []byte), in: make(chan []byte)},
		Link{out: make(chan []byte), in: make(chan []byte)}}
	l.remote.conn, err = startConn(rS, rT)
	misc.CheckError(err, "Establishing remote connection", true)
	go handleConn(&l.remote, &l.local, DST_R)
	if len(lS) != 0 {
		l.local.conn, err = startConn(lS, lT)
		misc.CheckError(err, "Establishing local connection", true)
		go handleConn(&l.local, &l.remote, DST_RL)
	}
	return &l
}

func startConn(s string, t bool) (*net.TCPConn, error) {

	log.Print("Starting connection on ", s)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", s)
	misc.CheckError(err, "Starting connection in Linker", true)
	if t {
		//TODO for accept
		log.Print("Creating server/Listening for client")
		listener, err := net.ListenTCP("tcp", tcpAddr)
		misc.CheckError(err, "Creating listener for Server in Linker", true)
		log.Print("Client acepted")
		return listener.AcceptTCP()
	}
	log.Print("Dialing...")
	return net.DialTCP("tcp", nil, tcpAddr)
}

func (l *Linker) Send(b []byte) {

	r := b[0]
	if r&(DST_R|DST_RL) > 0 {
		l.remote.out <- b
	}
	if r&DST_L > 0 {
		l.local.out <- b
	}
}

func (l *Linker) RegisterChannel(r bool) chan []byte {

	if r {
		return l.remote.in
	}
	return l.local.in
}

func handleConn(l, o *Link, t uint8) {

	var quit chan byte

	defer l.conn.Close()
	defer close(l.in)
	defer close(l.out)

	defer func() { quit <- 1 }()

	buff := make([]byte, 32)
	go func() {
		for {
			select {
			case <-quit:
				return
			case b := <-l.out:
				_, err := l.conn.Write(append([]byte{MAGIC}, b...))
				if misc.CheckError(err, "Sending data to conn", false) != nil {
					return
				}
			}
		}
	}()
	for {
		_, err := l.conn.Read(buff[0:])
		if misc.CheckError(err, "Receiving data from conn", false) != nil {
			return
		}
		if buff[0] != MAGIC || len(buff) < 3 {
			fmt.Println("Wrong packet")
			return
		}
		//TODO redirect to local or remote if necessary
		/*
			switch m := buff[1]; {
			case m &^ t > 0 && o != nil:
				fmt.Println("other")
				o.in <- buff[0:]
			default:
		*/
		l.in <- buff[3:]
		/*
			}
		*/
	}
}

func (l *Linker) Stop()  {
	if l.remote.conn != nil {
		l.remote.conn.Close()
		close(l.remote.in)
		close(l.remote.out)
	}
	if l.local.conn != nil {
		l.local.conn.Close()
		close(l.local.in)
		close(l.local.out)
	}
}
