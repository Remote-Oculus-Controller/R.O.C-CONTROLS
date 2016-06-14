package roc

import (
	"fmt"
	"github.com/Happykat/R.O.C-CONTROLS/misc"
)

func (roc *Roc) forward(data []byte) error {

	v, err := misc.DecodeInt(data)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("forward !!!!", v)
	return nil
}

func (roc *Roc) stop() byte {
	fmt.Println("stop !!!!")
	return 200
}
