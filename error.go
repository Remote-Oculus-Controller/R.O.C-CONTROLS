package roc

import (
	"github.com/Remote-Oculus-Controller/proto"
	"github.com/Remote-Oculus-Controller/proto/go"
)

//NewError create a Packet struct conform to the protocol that indicate an error occurred
func NewError(code rocproto.Error_Codes, msg string, dst int32) *rocproto.Packet {

	p := goPack.Prepare(uint32(code), rocproto.Packet_ERROR, rocproto.Packet_CONTROL_SERVER, rocproto.Packet_Section(dst))
	p.Err = &rocproto.Error{Code: uint32(code), Msg: msg}
	return p
}
