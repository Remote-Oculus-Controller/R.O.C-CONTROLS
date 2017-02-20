package network

import (
	"github.com/Remote-Oculus-Controller/R.O.C-CONTROLS/misc"
	"github.com/Remote-Oculus-Controller/proto"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

//WSSrv websocket server implementation for roc network interface
type WSSrv struct {
	ws *websocket.Conn
	*RocNet
}

//NewWSSrv return a pointer to a Websocket server
func NewWSSrv(r *RocNet) *WSSrv {
	return &WSSrv{RocNet: r}
}

//Start a websocket connection
//A websocket http endpoint is created on localhost:port/controls
func (ws *WSSrv) Start() {

	log.Printf("Starting websocket on %v/controls\n", ws.ip)
	http.HandleFunc("/controls", ws.listen)
	err := http.ListenAndServe(ws.ip, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//Stop and free
func (ws *WSSrv) Stop() {

	ws.ws.Close()
	ws.ws = nil
	ws.open = false
}

func (ws *WSSrv) listen(w http.ResponseWriter, r *http.Request) {

	if ws.ws != nil {
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
	ws.ws = c
	ws.open = true
	go ws.handle(rocproto.Packet_CONTROL_SERVER | rocproto.Packet_VIDEO_CLIENT)
	return
}

//Handle websocket connection
//
//The t parameters contain the section to which the handle is associated with Controls<=>Server Video / Controls<=>Client Videos
func (ws *WSSrv) handle(t rocproto.Packet_Section) {

	defer ws.Stop()

	quit := make(chan bool)
	go ws.receive(quit)
	ws.send(quit)
}

func (ws *WSSrv) receive(quit chan bool) {

	defer func() { quit <- true }()
	for {
		_, buff, err := ws.ws.ReadMessage()
		if misc.CheckError(err, "Receiving data from conn", false) != nil {
			return
		}
		ws.orderPacket(buff)
	}
}

func (ws *WSSrv) send(quit chan bool) {

	for {
		select {
		case <-quit:
			return
		case m := <-ws.out:
			b, err := proto.Marshal(m)
			if misc.CheckError(err, "linker.go/handleWS", false) != nil {
				if m == nil {
					return
				}
				continue
			}
			err = ws.ws.WriteMessage(websocket.BinaryMessage, b)
			if misc.CheckError(err, "linker.go/handleWS", false) != nil {
				return
			}
		}
	}
}

//Connected ...
func (ws *WSSrv) Connected() bool {

	if ws.ws != nil {
		return true
	}
	return false
}
