package roc

import (
	"fmt"
)

func (roc *Roc) forward(spd byte) byte {
	spd = spd * 2
	fmt.Println("forward !!!!")
	return 200
}

func (roc *Roc) stop() byte {
	fmt.Println("stop !!!!")
	return 200
}
