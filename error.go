package roc

import (
	"github.com/Remote-Oculus-Controller/proto"
	"github.com/Remote-Oculus-Controller/proto/go"
)

func NewError(code int, msg string) *rocproto.Packet {

	p := goPack.Prepare(code, rocproto.Packet_ERROR, rocproto.Packet_CONTROL_SERVER, rocproto.Packet_VIDEO_SERVER)
	p.Err = &rocproto.Error{Code: code, Msg: msg}
	return p
}
