package misc

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
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

func DecodeBytes(key []byte, v interface{}) (interface{}, error){

	var buff bytes.Buffer

	fmt.Println(key)
	buff.Write(key)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&v)
	if err != nil {
		fmt.Println("error")
		return nil, err
	}
	return v, nil
}

func DecodeInt(key []byte, v int) (int, error) {

	var buff bytes.Buffer

	buff.Write(key)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&v)
	if err != nil {
		fmt.Println("error")
		return 0, err
	}
	return v, nil
}

func DecodeFloat32(key []byte) (float32, error) {

	var buff bytes.Buffer
	var r	float32

	buff.Write(key)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&r)
	if err != nil {
		log.Println("Error decoding float32 from bytes")
		return 0, err
	}
	return r, nil
}