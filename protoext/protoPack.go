package protoext

import "github.com/Happykat/R.O.C-CONTROLS"

const ()

func Prepare(id uint32, t_ roc.Packet_Type, s_, d_ roc.Packet_Section) (p *roc.Packet) {

	t := uint32(t_)
	s := uint32(s_)
	d := uint32(d_)

	p = new(roc.Packet)
	p.Magic = roc.MAGIC
	p.Header = t<<roc.Shift_TYPE | s<<roc.Shift_SEND | d
	p.ID = id
	return p
}

func ReverseTo(p *roc.Packet, t roc.Packet_Type) *roc.Packet {

	s := uint32(p.Header) & roc.Mask_SEND >> roc.Shift_SEND
	d := uint32(p.Header) & roc.Mask_DST << roc.Shift_SEND
	p.Header = uint32(t)<<roc.Shift_TYPE | s | d
	return p
}

func GetSender(p *roc.Packet) string {

	s := uint32(p.Header) & roc.Mask_SEND >> roc.Shift_SEND
	return roc.Packet_Section_name[int32(s)]
}
