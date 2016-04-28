package misc

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

func EncodeBytes(key interface{}) ([]byte, error) {

	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodeBytes(key []byte) ([]byte, error){

	var buf bytes.Buffer

	fmt.Println(key)
	dec := gob.NewDecoder(&buf)
	err := dec.Decode(&key)
	if err != nil {
		fmt.Println("error")
		return nil, err
	}
	return buf.Bytes(), nil
}