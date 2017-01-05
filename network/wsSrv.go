package network

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"github.com/Remote-Oculus-Controller/proto"
	"github.com/Remote-Oculus-Controller/R.O.C-CONTROLS/misc"
	"github.com/golang/protobuf/proto"
)

type WsSrv struct {
	ws *websocket.Conn
	*RocNet
}

func NewWsSrv(r *RocNet) *WsSrv{
	return &WsSrv{RocNet: r}
}

//Start a websocket connection
//A websocket http endpoint is created on localhost:port/controls
func (l *WsSrv) Start() {

	log.Printf("Starting websocket on %v/controls\n", l.ip)
	http.HandleFunc("/controls", l.listen)
	err := http.ListenAndServe(l.ip, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (l *WsSrv) Stop() {

	l.ws.Close()
	l.ws = nil
	l.open = false
}

func (l *WsSrv) listen(w http.ResponseWriter, r *http.Request) {

	if l.ws != nil {
		er := "Remote connection already taken !!"
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
	l.open = true
	go l.handle(rocproto.Packet_CONTROL_SERVER|rocproto.Packet_VIDEO_CLIENT)
	return
}

//Handle websocket connection
//
//The t parameters contain the section to which the handle is associated with Controls<=>Server Video / Controls<=>Client Videos
func (l *WsSrv) handle(t rocproto.Packet_Section) {

	defer l.Stop()

	quit := make(chan bool)
	go l.receive(quit)
	l.send(quit)
}

func (s *WsSrv) receive(quit chan bool)  {

	defer func() { quit <- true }()
	for {
		_, buff, err := s.ws.ReadMessage()
		if misc.CheckError(err, "Receiving data from conn", false) != nil {
			return
		}
		s.orderPacket(buff)
	}
}

func (s *WsSrv) send(quit chan bool) {

	for {
		select {
		case <-quit:
			return
		case m := <-s.out:
			log.Printf("Sending ==>	%v\n", m)
			b, err := proto.Marshal(m)
			if misc.CheckError(err, "linker.go/handleWS", false) != nil {
				if m == nil {
					return
				}
				continue
			}
			err = s.ws.WriteMessage(websocket.BinaryMessage, b)
			if misc.CheckError(err, "linker.go/handleWS", false) != nil {
				return
			}
		}
	}
}

func (ws *WsSrv) Connected() bool{

	if (ws.ws != nil) {
		return true
	}
	return false
}