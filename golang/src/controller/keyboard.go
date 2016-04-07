package controller

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/keyboard"
)

type Keyboard struct {
	Controller
}

const (
	KEYBOARD_CF = "keyboard_map.json"
)

func (k Keyboard) Type() string {
	return "Keyboard"
}

func NewKeyboard(out chan byte) *Keyboard {

	k := Keyboard{}
	k.out = out
	k.mapControl(KEYBOARD_CF)
	keys := keyboard.NewKeyboardDriver("keyboard")
	fmt.Println(k.cmap)
	work := func() {
		gobot.On(keys.Event("key"), func(data interface{}) {
			key := data.(keyboard.KeyEvent)
			k.packet([]byte{(k.cmap[string(key.Key)])})
		})
	}

	k.robot = gobot.NewRobot("keyboard",
		[]gobot.Connection{},
		[]gobot.Device{keys},
		work,
	)
	return &k
}
