package controller

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/keyboard"
	"github.com/Happykat/R.O.C-CONTROLS"
)

type Keyboard struct {
	roc.Controller
}

var KEYBOARD_CF = roc.CF_DIR + "keyboard_map.json"

func (k Keyboard) Type() string {
	return "Keyboard"
}

func NewKeyboard() *Keyboard {

	k := new(Keyboard)
	k.MapControl(KEYBOARD_CF)
	keys := keyboard.NewKeyboardDriver("keyboard")
	work := func() {
		gobot.On(keys.Event("key"), func(data interface{}) {
			key := data.(keyboard.KeyEvent)
			k.Packet(string(key.Key), nil)
		})
	}
	k.Robot = gobot.NewRobot("keyboard",
		[]gobot.Connection{},
		[]gobot.Device{keys},
		work,
	)
	return k
}
