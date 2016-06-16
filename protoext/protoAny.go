package protoext

import (
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

func UnpackAny(b *any.Any, m proto.Message) error {

	err := proto.Unmarshal(b.Value, m)
	return err
}
