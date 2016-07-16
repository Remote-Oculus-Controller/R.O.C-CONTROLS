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

	if b == nil || m == nil {
		fmt.Printf("Cannot unpack message, one paramters is null any %v message %v", b, m)
		return nil
	}
	err := proto.Unmarshal(b.Value, m)
	return err
}
