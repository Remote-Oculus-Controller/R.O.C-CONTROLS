package misc

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

func DecodeInt(key []byte) (int, error) {

	var buff bytes.Buffer
	var v int

	buff.Write(key)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&v)
	if err != nil {
		fmt.Println("error")
		return 0, err
	}
	return v, nil
}
