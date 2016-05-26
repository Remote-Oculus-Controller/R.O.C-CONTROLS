package roc

import (
	"github.com/hybridgroup/gobot"
	"log"
	"fmt"
)

type RocClient struct {
	*gobot.Gobot
	l *Linker
}

func NewClient(lS, rS string, lT, rT bool) *RocClient {

	client := new(RocClient)
	client.Gobot = gobot.NewGobot()
	client.l = NewLinker(lS, rS, lT, rT)
	go func() {
		for {
			select {
			case b := <-client.l.remote.in:
				fmt.Println(b)
			}
		}
	}()
	return client
}

func (roc *RocClient) Start() error {

	roc.l.Start()
	errs := roc.Gobot.Start()
	if (errs != nil) {
		for _, err := range errs {
			log.Println(err)
		}
		panic(fmt.Sprintln("Panic starting RocClient"))
	}
	return nil
}

func (client *RocClient) SetController(c *Controller) {

	c.send=client.l.Send
	client.AddRobot(c.Robot)
}