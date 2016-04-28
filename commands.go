package roc

import (
	"fmt"
)

func (roc *Roc) forward(data []byte) (byte, error){
	fmt.Println("forward !!!!", data)
	return 200, nil
}

func (roc *Roc) stop() byte {
	fmt.Println("stop !!!!")
	return 200
}
