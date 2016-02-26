package linker

import (
	"fmt"
	"net"
)

type Linker struct {
	t, ip, port string
	direction   bool
	in, out     chan byte
	conn        net.Conn
}

func NewLinker(_t, _ip, _port string, _d bool) *Linker {
	linker := new(Linker)
	linker.t = _t
	linker.ip = _ip
	linker.port = _port
	linker.direction = _d
	return linker
}

func (l *Linker) Start() error {

	var err error

	fmt.Println(l.t + ": Starting connection...")
	if l.direction == true {
		ln, err := net.Listen("tcp", l.ip+":"+l.port)
		if err != nil {
			fmt.Println(l.t + " Error: " + err.Error())
			return err
		}
		fmt.Println(l.t + ": Socket established")
		l.conn, err = ln.Accept()
		if err != nil {
			fmt.Println(l.t + " Error: " + err.Error())
			return err
		}
		fmt.Println(l.t + ": Client accepted")
	} else {
		fmt.Println(l.t + ": Dialing server at adress " + l.ip + ":" + l.port)
		l.conn, err = net.Dial("tcp", l.ip+":"+l.port)
		if err != nil {
			fmt.Println(l.t + " Error: " + err.Error())
			return err
		}
		fmt.Println(l.t + ": Connected")
	}
	l.in = make(chan byte)
	l.out = make(chan byte)
	handleConnection(l.conn, l.in, l.out)
	return nil
}

func (l *Linker) Stop() error {
	fmt.Println(l.t + ": Closing connection...")
	if l.conn != nil {
		err := l.conn.Close()
		if err != nil {
			fmt.Println(l.t + ": Error during link disconnection")
			return err
		}
		fmt.Println(l.t + ": Connection closed")
	} else {
		fmt.Println(": Connection already closed or not started.")
	}
	return nil
}
