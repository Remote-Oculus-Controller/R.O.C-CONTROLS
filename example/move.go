package main

import (
	"github.com/Happykat/R.O.C-CONTROLS"
	"github.com/Happykat/R.O.C-CONTROLS/protoext"
	"github.com/Happykat/R.O.C-CONTROLS/robots"
	"log"
)

func main() {

	var err error

	m := robots.NewMotion()
	go m.Start()

	p := robots.Gyro{X: 60, Y: 100}
	r := roc.Packet{}
	r.Payload, err = protoext.PackAny(p)
	if err != nil {
		log.Fatal(err.Error())
	}
}
