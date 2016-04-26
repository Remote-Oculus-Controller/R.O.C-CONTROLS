package controller

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/api"
	"linker"
	"log"
	"misc"
	"parser"
)

type Controller struct {
	gbot  *gobot.Gobot
	robot *gobot.Robot
	cmap  map[string]parser.Cmd
	link  *linker.Linker
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

//TODO see to buffer
func (c *Controller) packet(code byte, data interface{}) {

	b, err := misc.GetBytes(data)
	if err != nil {
		panic(fmt.Sprintf(err.Error()))
	}
	b = append([]byte{linker.CMD | linker.DST_R | linker.DST_L | linker.MV, byte(len(b) + 1), code}, b...)
	c.link.Send(b)
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

	c.cmap = make(map[string]parser.Cmd)
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
			c.cmap[k] = vv
		} else {
			return ControllerError{"parseControl", v.(string) + " key can't be found in command.json", nil}
		}
	}
	log.Println(c.cmap)
	return nil
}
