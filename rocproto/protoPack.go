package rocproto

func Prepare(id uint32, t_ Packet_Type, s_, d_ Packet_Section) (p *Packet) {

	t := uint32(t_)
	s := uint32(s_)
	d := uint32(d_)

	p = new(Packet)
	p.Magic = Packet_MAGIC_Number
	p.Header = t<<Packet_SHIFT | s<<Packet_SHIFT_SENT | d
	p.ID = id
	return p
}

func ReverseTo(p *Packet, t Packet_Type) *Packet {

	s := uint32(p.Header) & Packet_MASK_SENT >> Packet_SHIFT_SENT
	d := uint32(p.Header) & Packet_MASK_DEST << Packet_SHIFT_SENT
	p.Header = uint32(t)<<Packet_SHIFT | s | d
	return p
}

func GetSender(p *Packet) string {

	s := uint32(p.Header) & Packet_MASK_SENT >> Packet_SHIFT_SENT
	return Packet_Section_name[int32(s)]
}
