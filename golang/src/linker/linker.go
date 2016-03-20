package linker

import (
	"log"
	"net"
)

type Linker struct {
	t, ip, port string
	direction   bool
	In, Out     chan byte
	conn        net.Conn
}

func NewLinker(_t, _ip, _port string, _d bool) *Linker {
	linker := new(Linker)
	linker.t = _t
	linker.ip = _ip
	linker.port = _port
	linker.direction = _d
	linker.In = make(chan byte)
	linker.Out = make(chan byte)
	return linker
}

func (l *Linker) Start() error {

	var err error

	log.Println(l.t + ": Starting connection...")
	if l.direction == true {
		ln, err := net.Listen("tcp", l.ip+":"+l.port)
		if err != nil {
			log.Println(l.t + " Error: " + err.Error())
			return err
		}
		log.Println(l.t + ": Socket established")
		l.conn, err = ln.Accept()
		if err != nil {
			log.Println(l.t + " Error: " + err.Error())
			return err
		}
		log.Println(l.t + ": Client accepted")
	} else {
		log.Println(l.t + ": Dialing server at adress " + l.ip + ":" + l.port)
		l.conn, err = net.Dial("tcp", l.ip+":"+l.port)
		if err != nil {
			log.Println(l.t + " Error: " + err.Error())
			return err
		}
		log.Println(l.t + ": Connected")
	}
	handleConnection(l.conn, l.In, l.Out)
	return nil
}

func (l *Linker) Stop() error {
	log.Println(l.t + ": Closing connection...")
	if l.conn != nil {
		err := l.conn.Close()
		if err != nil {
			log.Println(l.t + ": Error during link disconnection")
			return err
		}
		log.Println(l.t + ": Connection closed")
	} else {
		log.Println(l.t + ": Connection already closed or not started.")
	}
	return nil
}
