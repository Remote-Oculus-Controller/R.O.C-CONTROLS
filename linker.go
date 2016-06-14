package roc

import (
	"errors"
	"github.com/Happykat/R.O.C-CONTROLS/misc"
	"log"
	"net"
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
	lIp, rIp      string
	lT, rT        bool
}

type Link struct {
	conn    *net.TCPConn
	out, in chan []byte
}

func NewLinker(lS, rS string, lT, rT bool) *Linker {

	l := Linker{local: Link{out: make(chan []byte, 100), in: make(chan []byte, 100)},
		remote: Link{out: make(chan []byte, 100), in: make(chan []byte, 100)},
		lIp:    lS, lT: lT, rIp: rS, rT: rT}
	return &l
}

func (l *Linker) Start() {
	if l.lIp != "" {
		log.Print("Staring local work")
		go l.local.startConn(l.lIp, l.lT, &l.remote, DST_RL)
	}
	log.Println("Starting remote work")
	go l.remote.startConn(l.rIp, l.rT, &l.local, DST_R)
}

//TODO timeout connection and try
func (l *Link) startConn(s string, m bool, o *Link, t uint8) {

	defer close(l.in)
	defer close(l.out)

	log.Print("Starting connection on ", s)
	tcpAddr, err := net.ResolveTCPAddr("tcp", s)
	misc.CheckError(err, "resolving address in linker.go/startConn", true)
	if m {
		log.Println("Listening on", tcpAddr.String())
		listener, err := net.ListenTCP("tcp", tcpAddr)
		misc.CheckError(err, "listening in linker.go/startConn", true)
		for {
			log.Println("Looking for a client...")
			l.conn, err = listener.AcceptTCP()
			misc.CheckError(err, "Accepting client in linker.go/startCnn", true)
			log.Print("Connection acepted")
			l.handleConn(o, t)
			log.Println("Closing connection")
			l.conn.Close()
		}
	} else {
		for {
			log.Print("Dialing...")
			l.conn, err = net.DialTCP("tcp", nil, tcpAddr)
			misc.CheckError(err, "Dialing adresse in linker.go/startConn", true)
			l.handleConn(o, t)
			l.conn.Close()
		}
	}
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

func (l *Link) handleConn(o *Link, t uint8) {

	var quit chan byte

	buff := make([]byte, 128)
	go func() {

		defer func() { quit <- 1 }()
		for {
			_, err := l.conn.Read(buff[0:])
			if misc.CheckError(err, "Receiving data from conn", false) != nil {
				return
			}
			if buff[0] != MAGIC || len(buff) < 3 {
				log.Println("Wrong packet")
				continue
			}
			if buff[1]&t != 0 {
				l.in <- buff[3:]
			}
			if (buff[1]&DST_ALL)&^(t|DST_L) != 0 && o != nil {
				o.out <- buff[0:]
			}
		}
	}()
	for {
		select {
		case <-quit:
			return
		case b := <-l.out:
			_, err := l.conn.Write(append([]byte{MAGIC}, b...))
			if misc.CheckError(err, "linker.go/handleConn", false) != nil {
				return
			}
		}
	}
}

func (l *Linker) Stop() {
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
