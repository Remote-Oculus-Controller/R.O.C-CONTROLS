package controller

import (
	"fmt"
	"github.com/Happykat/R.O.C-CONTROLS"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/keyboard"
)

type Keyboard struct {
	Controller
}

var KEYBOARD_CF = CF_DIR + "keyboard_map.json"

func (k Keyboard) Type() string {
	return "Keyboard"
}

func NewKeyboard(link *roc.Linker) *Keyboard {

	k := new(Keyboard)
	k.link = link
	k.mapControl(KEYBOARD_CF)
	fmt.Println(k.cmap)
	keys := keyboard.NewKeyboardDriver("keyboard")
	work := func() {
		gobot.On(keys.Event("key"), func(data interface{}) {
			key := data.(keyboard.KeyEvent)
			p := k.cmap[string(key.Key)]
			k.packet(p.Code, p.Default)
		})
	}
	k.robot = gobot.NewRobot("keyboard",
		[]gobot.Connection{},
		[]gobot.Device{keys},
		work,
	)
	return k
}
