package roc

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
)

func PackAny(m proto.Message) (*any.Any, error) {

	a := new(any.Any)
	a.TypeUrl = proto.MessageName(m)
	p, err := proto.Marshal(m)
	if err != nil {
		return nil, err
	}
	a.Value = p
	return a, nil
}

//TODO check message url
func UnpackAny(b *any.Any, m proto.Message) error {

	var err error

	err = nil
	if b != nil && m != nil {
		err = proto.Unmarshal(b.Value, m)
	} else {
		fmt.Printf("Cannot parse proto.Message, one memeber is null")
	}
	return err
}
