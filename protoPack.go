package roc

func Prepare(id uint32, t_ Packet_Type, s_, d_ Packet_Section) (p *Packet) {

	t := uint32(t_)
	s := uint32(s_)
	d := uint32(d_)

	p = new(Packet)
	p.Magic = MAGIC
	p.Header = t<<Shift_TYPE | s<<Shift_SEND | d
	p.ID = id
	return p
}

func ReverseTo(p *Packet, t Packet_Type) *Packet {

	s := uint32(p.Header) & Mask_SEND >> Shift_SEND
	d := uint32(p.Header) & Mask_DST << Shift_SEND
	p.Header = uint32(t)<<Shift_TYPE | s | d
	return p
}

func GetSender(p *Packet) string {

	s := uint32(p.Header) & Mask_SEND >> Shift_SEND
	return Packet_Section_name[int32(s)]
}
