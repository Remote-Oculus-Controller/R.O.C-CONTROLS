package controller

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/api"
	"log"
	"parser"
)

type Controller struct {
	robot *gobot.Robot
	gbot  *gobot.Gobot
	cmap  map[string]byte
}

type Bind struct {
	name string
	cmd  byte
}

func (c *Controller) Type() string {
	return "Controller"
}

//TODO redo error
func (c *Controller) Start() error {

	if c.robot != nil {
		c.gbot = gobot.NewGobot()
		a := api.NewAPI(c.gbot)
		a.Debug()
		a.Start()
		c.gbot.AddRobot(c.robot)
		err := c.gbot.Start()
		return err[0]
	}
	log.Println("Pls init controller before starting")
	return nil
}

func (c *Controller) Stop() error {

	c.gbot.Stop()
	return nil
}

// TODO change association to direct function, no map of pointer to function !
func (c *Controller) parseControl(fp string) error {

	c.cmap = make(map[string]byte)
	b, err := parser.Decode(fp)
	if err != nil {
		log.Println("Failed to parse", fp)
		return err
	}
	m, err := parser.RobotCommand("command.json")
	if err != nil {
		log.Println("Failed to parse command.json.")
		return err
	}
	for k, v := range b {
		log.Println("pass", k, v)
		if vv, ok := m[v.(string)]; ok {
			c.cmap[k] = vv.Code
		} else {
			return ControllerError{"parseControl", v.(string) + " key can't be found in command.json", nil}
		}
	}
	log.Println(c.cmap)
	return nil
}
