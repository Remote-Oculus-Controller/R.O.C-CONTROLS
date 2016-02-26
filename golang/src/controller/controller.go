package controller

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/api"
	"parser"
)

type Controller struct {
	robot *gobot.Robot
	gbot  *gobot.Gobot
	cmap  map[string]parser.Cmd
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
	fmt.Println("Pls init controller before starting")
	return nil
}

func (c *Controller) Stop() error {

	c.gbot.Stop()
	return nil
}

func (c *Controller) parseControl(fp string) error {

	c.cmap = make(map[string]parser.Cmd)
	b, err := parser.Decode(fp)
	if err != nil {
		fmt.Println("Failed to parse", fp)
		return err
	}
	m, err := parser.RobotCommand("command.json")
	if err != nil {
		fmt.Println("Failed to parse command.json.")
		return err
	}
	for k, v := range b {
		fmt.Println("pass", k, v)
		if vv, ok := m[v.(string)]; ok {
			fmt.Println("key found", vv, ok)
			c.cmap[k] = vv
		} else {
			return ControllerError{"parseControl", v.(string) + " key can't be found in command.json", nil}
		}
	}
	return nil
}
