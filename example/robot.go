package main

import (
	"flag"

	"github.com/Happykat/R.O.C-CONTROLS"
	"github.com/Happykat/R.O.C-CONTROLS/robots"
)

var remote string
var local string
var remote_t bool
var local_t bool

func main() {
	flag.Parse()
	r := roc.NewRoc(local, remote, local_t, remote_t)
	r.AddRocRobot(robots.NewGPS().RocRobot)
	r.AddRocRobot(robots.NewMotion().RocRobot)
	r.AddRocRobot(robots.NewAI(r).RocRobot)
	r.Start()
}

func init() {
	flag.StringVar(&remote, "r", ":8001", "ip:port parameter, ip:port")
	flag.StringVar(&local, "l", "", "ip:port parameter, ip:port")
	flag.BoolVar(&remote_t, "rT", true,
		"Set this side of connection as a server(true)or client(false)")
	flag.BoolVar(&local_t, "lT", false,
		"Set this side of connection as a server(true)or client(false)")
}
