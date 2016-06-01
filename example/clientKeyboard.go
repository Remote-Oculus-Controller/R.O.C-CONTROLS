package main

import (
"flag"
"github.com/Happykat/R.O.C-CONTROLS"
"github.com/Happykat/R.O.C-CONTROLS/controller"
)

var remote	string
var local	string
var remote_t	bool
var local_t	bool

func main() {

	flag.Parse()
	client := roc.NewClient(local, remote, local_t, remote_t)
	client.SetController(&controller.NewKeyboard().Controller)
	client.Start()
}

func init() {
	flag.StringVar(&remote, "r", "127.0.0.1:4343", "ip:port parameter, ip:port")
	flag.StringVar(&local, "l", "", "ip:port parameter, ip:port")
	flag.BoolVar(&remote_t, "rT", false,
	"Set this side of connection as a server(true)or client(false)")
	flag.BoolVar(&local_t, "lT", false,
	"Set this side of connection as a server(true)or client(false)")
}
