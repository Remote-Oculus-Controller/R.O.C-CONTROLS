package roc

import (
	"fmt"
)

func (roc *Roc) cmd(b []byte) {
	switch b[0] {
	case 0x0A:
		switch b[1] {
		case 0x01:
			roc.forward(255)
		case 0x00:
			roc.stop()
		}
	}
}
func (roc *Roc) forward(spd byte) byte {
	spd = spd * 2
	fmt.Println("forward !!!!")
	return 200
}

func (roc *Roc) stop() byte {
	fmt.Println("stop !!!!")
	return 200
}
