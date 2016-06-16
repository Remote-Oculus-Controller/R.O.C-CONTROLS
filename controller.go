package roc

import (
	/*
		"fmt"
	*/
	"github.com/hybridgroup/gobot"
	"log"
	"os"
	/*
		"github.com/Happykat/R.O.C-CONTROLS/misc"
	*/
	"errors"
)

type Controller struct {
	*gobot.Robot
	cmap map[string]Cmd
	send func(*Packet) error
}

var CF_DIR = func() string {
	dir, err := os.Getwd()
	if err != nil {
		return "error in work directory"
	}
	return dir
}() + "/config/"
var CMD_FILE = CF_DIR + "command.json"

func (c *Controller) Type() string {
	return "Controller"
}

func (c *Controller) Packet(event string, data interface{}) {

	/*
		var b []byte
		var err error

		e, ok := c.cmap[event]
		if (!ok) {
			return
		}
		_, def := data.(bool)
		if (data == nil || def) {
			b, err = misc.EncodeBytes(e.Default)
		} else {
			b, err = misc.EncodeBytes(data)
		}
		if err != nil {
			panic(fmt.Sprintf(err.Error()))
		}
	*/
	/*
		//TODO redo controller
		b = append([]byte{CMD | DST_R | MV, e.Code}, b...)
	*/
	c.send(&Packet{})
}

func (c *Controller) MapControl(file string) error {
	log.Println("Mapping", c.Type(), "for robot control\nStart parsing :", file)
	err := c.parseControl(file)
	if err != nil {
		log.Println("Failed to parse control.")
		return err
	}
	return nil
}

func (c *Controller) parseControl(fp string) error {

	c.cmap = make(map[string]Cmd)
	b, err := DecodeJsonFile(fp)
	if err != nil {
		log.Println("Failed to parse", fp)
		return err
	}
	m, err := ParseCommands(CMD_FILE)
	if err != nil {
		log.Println("Failed to parse command.json.")
		return err
	}
	for k, v := range b {
		if vv, ok := m[v.(string)]; ok {
			c.cmap[k] = vv
		} else {
			return errors.New("parseControl : " + v.(string) + " key can't be found in command.json")
		}
	}
	log.Println(c.cmap)
	return nil
}
