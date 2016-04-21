package linker

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestSend(t *testing.T) {

	serv := NewLinker("tcp", "127.0.0.1", "4343", true)
	client := NewLinker("tcp", "127.0.0.1", "4343", false)
	go serv.Start()
	fmt.Println("tt")
	time.Sleep(time.Second * 3)
	go client.Start()
	for client.conn != nil {

	}
	cases := []struct {
		in []byte
	}{
		{[]byte("Hello, world")},
		{[]byte("Hello, 世界")},
		{[]byte("")},
	}
	buff := new(bytes.Buffer)
	for _, c := range cases {
		fmt.Println("oo")
		client.Send(bytes.NewBuffer(c.in))
		l := <-serv.In
		for i := byte(0); i < l; i++ {
			buff.WriteByte(<-serv.In)
		}
		got := buff.Bytes()
		if bytes.Compare(c.in, got) == 0 {
			t.Errorf("Send(%q) == %q, want %q", c.in, got, c.in)
		}
	}
	serv.Stop()
	client.Stop()
}
