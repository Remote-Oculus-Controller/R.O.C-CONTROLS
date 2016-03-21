package controller

import (
	"github.com/hybridgroup/gobot"
	"log"
	"parser"
)

type Controller struct {
	gbot  *gobot.Gobot
	robot *gobot.Robot
	cmap  map[string]byte
	out   chan byte
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
		c.gbot.AddRobot(c.robot)
		err := c.gbot.Start()
		if err != nil {
			log.Println(err)
			return err[0]
		}
	} else {
		log.Println("Please use the New'controller' to create controller")
	}
	return nil
}

func (c *Controller) Stop() error {

	c.gbot.Stop()
	return nil
}

func (c *Controller) mapControl(file string) error {
	log.Println("Mapping", c.Type(), "for robot control\nStart parsing :", file)
	err := c.parseControl(file)
	if err != nil {
		log.Println("Failed to parse control.")
		return err
	}
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
		if vv, ok := m[v.(string)]; ok {
			c.cmap[k] = vv.Code
		} else {
			return ControllerError{"parseControl", v.(string) + " key can't be found in command.json", nil}
		}
	}
	log.Println(c.cmap)
	return nil
}
