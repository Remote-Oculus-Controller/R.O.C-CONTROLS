package roc

import (
	"fmt"
	"log"
	"github.com/hybridgroup/gobot"
	"os"
	"github.com/Happykat/R.O.C-CONTROLS/misc"
	"errors"
)

type Controller struct {
	*gobot.Robot
	cmap  map[string]Cmd
	send	func([]byte) int
}

var CF_DIR = os.Getenv("GOPATH") + "/src/github.com/Happykat/R.O.C-CONTROLS/config/"
var CMD_FILE = CF_DIR + "command.json"

func (c *Controller) Type() string {
	return "Controller"
}

//TODO see to buffer
func (c *Controller) Packet(event string, data interface{}) {

	var b []byte
	var err error

	e, ok := c.cmap[event]
	if (!ok) {
		return
	}
	code := e.Code
	_, def := data.(bool)
	if (data == nil || def) {
		fmt.Println("default", e.Default)
		b, err = misc.EncodeBytes(e.Default)
	} else {
		b, err = misc.EncodeBytes(data)
	}
	if err != nil {
		panic(fmt.Sprintf(err.Error()))
	}
	b = append([]byte{CMD | DST_R | MV, code}, b...)
	c.send(b)
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

// TODO change association to direct function, no map of pointer to function !
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