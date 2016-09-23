package roc

import (
	"testing"

	"time"

	"github.com/Remote-Oculus-Controller/R.O.C-CONTROLS/test"
	"github.com/Remote-Oculus-Controller/proto"
	"github.com/Remote-Oculus-Controller/proto/go"
	"github.com/golang/protobuf/proto"
)

type ConnParam struct {
	lIp, rIp string
	lT, rT   bool
}

var connParamDefault = ConnParam{
	lIp: "",
	lT:  false,
	rIp: ":8004",
	rT:  true,
}

var l *Linker = newLinker(connParamDefault.lIp, connParamDefault.rIp, connParamDefault.lT, connParamDefault.rT)

func TestLinker_Start(t *testing.T) {

	l.start()
	client := roc_test.ClientWs{}
	go client.Start(connParamDefault.rIp)
	<-time.After(time.Millisecond * 200)
	if l.local.conn != nil {
		t.Error("Local connection should not be set")
	}
	if l.remote.ws == nil {
		t.Error("Remote connection websocket should have been created")
	}
	client.Stop()
	<-time.After(time.Millisecond * 1200)
}

func TestLinker_SendRemote(t *testing.T) {

	client := roc_test.ClientWs{}
	go client.Start(connParamDefault.rIp)
	defer client.Stop()
	<-time.After(time.Millisecond * 200)

	p := goPack.Prepare(1, rocproto.Packet_COMMAND, rocproto.Packet_CONTROL_SERVER, rocproto.Packet_VIDEO_CLIENT)
	err := l.send(p)
	if err != nil {
		t.Error(err)
	}
	select {
	case <-time.After(time.Millisecond * 5):
		t.Error("Time out > 5ms")
	default:
		b := client.Rcv()
		p2 := rocproto.Packet{}
		err := proto.Unmarshal(b, &p2)
		if err != nil {
			t.Error(err)
		}
		if *p != p2 {
			t.Error("Message different from original")
		}
	}
}

func TestLinker_Stop(t *testing.T) {
	l.stop()
	in := <-l.remote.in
	out := <-l.local.in
	if l.remote.ws != nil && in != nil && out != nil {
		t.Error("All ressources have not been closed", l.remote)
	}
}
