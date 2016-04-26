package linker

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSend(t *testing.T) {

	serv := NewLinker("", "127.0.0.1:4343", false, true)
	ch := serv.RegisterChannel(true)
	client := NewLinker("", "127.0.0.1:4343", false, false)
	cases := []struct {
		in []byte
	}{
		{[]byte("Hello, world")},
		{[]byte("Hello, 世界")},
		{[]byte("")},
	}
	for _, c := range cases {
		fmt.Println("oo")
		client.Send(c.in)
		got := <-ch
		if bytes.Compare(c.in, got) == 0 {
			t.Errorf("Send(%q) == %q, want %q", c.in, got, c.in)
		}
	}
	serv.Stop()
	client.Stop()
}
