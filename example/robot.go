package main

import (
	"flag"

	"github.com/Remote-Oculus-Controller/R.O.C-CONTROLS"
	"github.com/Remote-Oculus-Controller/R.O.C-CONTROLS/robots"
)

var remote string
var local string
var remote_t bool
var local_t bool

func main() {
	flag.Parse()
	r := roc.NewRoc(local, remote, local_t, remote_t)
	r.AddRocRobot(robots.NewGPS().Robot)
	r.AddRocRobot(robots.NewMotion().Robot)
	r.AddRocRobot(robots.NewAI(r).Robot)
	r.Start()
}

func init() {
	flag.StringVar(&remote, "r", ":8001", "CLIENT ip:port parameter, ip:port")
	flag.BoolVar(&remote_t, "rT", true,
		"Set this side of connection as a server(true)or client(false)")
	flag.StringVar(&local, "l", "", "VIDEO SERVER ip:port parameter, ip:port")
	flag.BoolVar(&local_t, "lT", false,
		"Set this side of connection as a server(true)or client(false)")
}
