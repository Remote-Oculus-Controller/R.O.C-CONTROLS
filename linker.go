package roc

import (
	"errors"
	"log"
	"net"

	"net/http"

	"github.com/Happykat/R.O.C-CONTROLS/misc"
	"github.com/Happykat/R.O.C-CONTROLS/rocproto"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

type Linker struct {
	local, remote Link
	lIp, rIp      string
	lT, rT        bool
}

//TODO clear TCP
type Link struct {
	conn    *net.TCPConn
	ws      websocket.Conn
	out, in chan *rocproto.Packet
}

func NewLinker(lS, rS string, lT, rT bool) *Linker {

	l := Linker{local: Link{conn: nil, out: make(chan *rocproto.Packet, 100), in: make(chan *rocproto.Packet, 100)},
		remote: Link{conn: nil, out: make(chan *rocproto.Packet, 100), in: make(chan *rocproto.Packet, 100)},
		lIp:    lS, lT: lT, rIp: rS, rT: rT}
	return &l
}

func (l *Linker) Start() {
	if l.lIp != "" {
		log.Print("Staring local work")
		go l.local.startConnTCP(l.lIp, l.lT, &l.remote, rocproto.Packet_CONTROL_SERVER|rocproto.Packet_VIDEO_SERVER)
	}
	log.Println("Starting remote work")
	go l.remote.startConnWS(l.rIp, l.rT, &l.local, rocproto.Packet_CONTROL_SERVER|rocproto.Packet_VIDEO_CLIENT)
}

//TODO timeout connection and try
func (l *Link) startConnTCP(s string, m bool, o *Link, t rocproto.Packet_Section) {

	defer close(l.in)
	defer close(l.out)

	var listener *net.TCPListener

	log.Println("Starting connection on ", s)
	tcpAddr, err := net.ResolveTCPAddr("tcp", s)
	misc.CheckError(err, "resolving address in linker.go/startConn", true)
	for {
		if m {
			log.Println("Listening on", tcpAddr.String())
			listener, err = net.ListenTCP("tcp", tcpAddr)
			misc.CheckError(err, "listening in linker.go/startConn", true)
			log.Println("Looking for a client...")
			l.conn, err = listener.AcceptTCP()
			misc.CheckError(err, "Accepting client in linker.go/startCnn", true)
			log.Print("Connection acepted")
			listener.Close()
		} else {
			log.Print("Dialing...")
			l.conn, err = net.DialTCP("tcp", nil, tcpAddr)
			misc.CheckError(err, "Dialing adresse in linker.go/startConn", true)
		}
		l.handleConn(o, t)
		log.Println("Closing connection")
		listener.Close()
		l.conn.Close()
		l.conn = nil
	}
}

func (l *Link) startConnWS(s string, m bool, o *Link, t rocproto.Packet_Section) {

	defer close(l.in)
	defer close(l.out)

	log.Printf("Starting websocket on %v/controls", s)
	if m {
		http.HandleFunc("/controls", l.listenWS)
	}
}

func checkBuffer(r int, buff []byte) (m rocproto.Packet, err error) {

	if r > -1 {
		err = proto.Unmarshal(buff[0:r], m)
	} else {
		err = proto.Unmarshal(buff[0:], m)
	}
	if err != nil {
		log.Println("Cannot Unmarshall packet", err.Error())
		log.Println(r, buff[0:r])
		return nil, err
	}
	if m.Magic != MAGIC {
		log.Println("Wrong packet")
		return nil, errors.New("Wrong packet")
	}
	return m, nil
}

func routPacket(m rocproto.Packet, l, o *Link, t rocproto.Packet_Type) {
	if m.Header&uint32(t) != 0 {
		l.in <- m
	}
	if (m.Header&uint32(rocproto.Packet_MASK_DEST))&^uint32(t) != 0 && o.conn != nil {
		o.out <- m
	}
}

//TODO Insert buffer len and check
func (l *Link) handleConn(o *Link, t rocproto.Packet_Section) {

	buff := make([]byte, 128)
	quit := make(chan bool)
	go func() {

		defer func() { quit <- true }()

		m := new(rocproto.Packet)
		for {
			r, err := l.conn.Read(buff[0:])
			if misc.CheckError(err, "Receiving data from conn", false) != nil {
				return
			}
			if m, err = checkBuffer(r, buff); err != nil {
				continue
			}
			routPacket(m, l, o, t)
		}
	}()
	for {
		select {
		case <-quit:
			return
		case m := <-l.out:
			b, err := proto.Marshal(m)
			if misc.CheckError(err, "linker.go/handleConn", false) != nil {
				continue
			}
			_, err = l.conn.Write(b)
			if misc.CheckError(err, "linker.go/handleConn", false) != nil {
				return
			}
		}
	}
}

func (l *Link) listenWS(w http.ResponseWriter, r *http.Request) {

	if l.conn != nil {
		er := "Connection already taken !!"
		log.Println(er)
		w.Write([]byte(er))
		return
	}
	upgrader := websocket.Upgrader{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		er := "Error upgrading to websocket"
		log.Println(er)
		w.Write([]byte(er))
		return
	}
	log.Println("Client connected", r.RemoteAddr)
	l.ws = c
	return
}

func (l *Link) handleWS(o *Link, t rocproto.Packet_Section) {

	quit := make(chan bool)
	go func() {

		defer func() { quit <- true }()

		m := new(rocproto.Packet)
		for {
			_, buff, err := l.ws.ReadMessage()
			if misc.CheckError(err, "Receiving data from conn", false) != nil {
				return
			}
			if m, err = checkBuffer(-1, buff); err != nil {
				continue
			}
			routPacket(m, l, o, t)
		}
	}()
	for {
		select {
		case <-quit:
			return
		case m := <-l.out:
			b, err := proto.Marshal(m)
			if misc.CheckError(err, "linker.go/handleWS", false) != nil {
				continue
			}
			err = l.ws.WriteMessage(websocket.BinaryMessage, b)
			if misc.CheckError(err, "linker.go/handleWS", false) != nil {
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

func (l *Linker) Send(p *rocproto.Packet) error {

	if (p.Header&uint32(rocproto.Packet_VIDEO_CLIENT)) != 0 && l.remote.conn != nil {
		l.remote.out <- p
	}
	if (p.Header&uint32(rocproto.Packet_VIDEO_SERVER)) != 0 && l.local.conn != nil {
		if l.local.conn != nil {
			l.local.out <- p
		} else {
			return errors.New("Local connection not established could not send message")
		}
	}
	return nil
}

func (l *Linker) RegisterChannel(r bool) chan *rocproto.Packet {

	if r {
		return l.remote.in
	}
	return l.local.in
}
