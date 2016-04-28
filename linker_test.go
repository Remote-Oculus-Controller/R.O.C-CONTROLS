package roc

import (
	"testing"
	"time"
	"bytes"
)

/*
func TestStart(t *testing.T) {

	var serv *Linker

	go func() {
		serv = NewLinker("", "127.0.0.1:4343", false, true)
	}()
	time.Sleep(time.Second)
	client := NewLinker("", "127.0.0.1:4343", false, false)
	client.Stop()
	serv.Stop()
}
*/

func TestSendReceive(t *testing.T) {

	var serv *Linker
	var ch	 chan []byte

	go func() {
		serv = NewLinker("", "127.0.0.1:4343", false, true)
		ch = serv.RegisterChannel(true)
	}()

	time.Sleep(100*time.Millisecond)
	client := NewLinker("", "127.0.0.1:4343", false, false)
	cases := []struct {
		in []byte
	}{
		{[]byte{DST_R,0,0,0,0,0,0}},
		{[]byte{DST_R,1,1,1,1,1,1}},
		{[]byte{DST_R,255,255,255}},
	}
	for _, c := range cases {
		client.Send(c.in)
		got := <-ch
		if bytes.Compare(c.in, got) == 0 {
			t.Errorf("Send(%q) != %q", c.in, got)
		}
	}
	serv.Stop()
	client.Stop()
	time.Sleep(100*time.Millisecond)
}
