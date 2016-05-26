package main

import (
	"github.com/Happykat/R.O.C-CONTROLS"
	"flag"
	"github.com/Happykat/R.O.C-CONTROLS/robots"
)

var remote	string
var local	string
var remote_t	bool
var local_t	bool

func main() {
	flag.Parse()
	r := roc.NewRoc(local, remote, local_t, remote_t)
	//r.AddRobot(roc.NewMotion().Robot)
	r.AddRobot(robots.NewGPS().RocRobot)
	r.Start()
}

func init() {
	flag.StringVar(&remote, "r", "127.0.0.1:4343", "ip:port parameter, ip:port")
	flag.StringVar(&local, "l", "", "ip:port parameter, ip:port")
	flag.BoolVar(&remote_t, "rT", true,
		"Set this side of connection as a server(true)or client(false)")
	flag.BoolVar(&local_t, "lT", false,
		"Set this side of connection as a server(true)or client(false)")
}