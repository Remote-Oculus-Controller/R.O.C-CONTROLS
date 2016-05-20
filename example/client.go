package main

import (
	"github.com/Happykat/R.O.C-CONTROLS/controller"
	"github.com/Happykat/R.O.C-CONTROLS"
	"flag"
)

var remote	string
var local	string
var remote_t	bool
var local_t	bool

func main() {

	flag.Parse()
	robotL := roc.NewLinker(local, remote, local_t, remote_t)
	device := controller.NewKeyboard(robotL)
	//device := controller.NewDS3(robotL)
	//device := controller.NewXbox(robotL)
	//unityL := linker.NewLinker("unity", "127.0.0.1", "4343", true)
	robotL.Start()
	device.Start()
}

func init() {
	flag.StringVar(&remote, "r", "127.0.0.1:4343", "ip:port parameter, ip:port")
	flag.StringVar(&local, "l", "", "ip:port parameter, ip:port")
	flag.BoolVar(&remote_t, "rT", false,
		"Set this side of connection as a server(true)or client(false)")
	flag.BoolVar(&local_t, "lT", false,
		"Set this side of connection as a server(true)or client(false)")
}