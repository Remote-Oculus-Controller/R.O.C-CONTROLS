package roc

import (
	"fmt"
	"net"
	"log"
	"github.com/Happykat/R.O.C-CONTROLS/misc"
	"errors"
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

	l := Linker{Link{out: make(chan []byte), in: make(chan []byte)},
		Link{out: make(chan []byte), in: make(chan []byte)}}
	l.remote.conn = startConn(rS, rT)
	l.local.conn = nil
	if len(lS) != 0 {
		l.local.conn = startConn(lS, lT)
	}
	return &l
}

func (l *Linker) Start() {
	if (l.local.conn != nil) {
		log.Print("Staring local work")
		go handleConn(&l.local, &l.remote, DST_RL)
	}
	log.Println("Starting remote work")
	go handleConn(&l.remote, &l.local, DST_R)
}

func startConn(s string, t bool) (*net.TCPConn) {

	log.Print("Starting connection on ", s)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", s)
	misc.CheckError(err, "resolving address in linker.go/startConn", true)
	if t {
		//TODO for accept
		listener, err := net.ListenTCP("tcp", tcpAddr)
		misc.CheckError(err, "listening in linker.go/startConn", true)
		v, err := listener.AcceptTCP()
		misc.CheckError(err, "accepting client in linker.go/startConn", true)
		log.Print("Client acepted")
		return v
	}
	log.Print("Dialing...")
	v, err := net.DialTCP("tcp", nil, tcpAddr)
	misc.CheckError(err, "dialing adresse in linker.go/startConn", true)
	return v
}

func (l *Linker) Send(b []byte) error {

	r := b[0]
	if r&(DST_R|DST_RL) != 0 {
		l.remote.out <- b
	}
	if r&DST_L != 0 {
		if l.local.conn != nil {
			l.local.out <- b
		} else {
			return errors.New("Local connection not established could not send message")
		}
	}
	return nil
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

	buff := make([]byte, 128)
	go func() {
		for {
			select {
			case <-quit:
				return
			case b := <-l.out:
				_, err := l.conn.Write(append([]byte{MAGIC}, b...))
				if misc.CheckError(err, "Writing data on conn", false) != nil {
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
			continue
		}
		if (buff[1] & t != 0) {
			l.in <- buff[3:]
		}
		if ((buff[1] & DST_ALL) &^ (t | DST_L) != 0 && o != nil) {
			o.out <- buff[0:]
		}
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
