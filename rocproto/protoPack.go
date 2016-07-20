package rocproto

const (
	MAGIC      = uint32(Packet_MAGIC_Number)
	SHIFT      = uint32(Packet_SHIFT)
	SHIFT_SEND = uint32(Packet_SHIFT_SEND)
	MASK_SEND  = uint32(Packet_MASK_SEND)
	MASK_DEST  = uint32(Packet_MASK_DEST)
)

func Prepare(id uint32, t_ Packet_Type, s_, d_ Packet_Section) (p *Packet) {

	t := uint32(t_)
	s := uint32(s_)
	d := uint32(d_)

	p = new(Packet)
	p.Magic = MAGIC
	p.Header = t<<SHIFT | s<<SHIFT_SEND | d
	p.ID = id
	return p
}

func ReverseTo(p *Packet, t Packet_Type) *Packet {

	s := uint32(p.Header) & MASK_SEND >> SHIFT_SEND
	d := uint32(p.Header) & MASK_DEST << SHIFT_SEND
	p.Header = uint32(t)<<SHIFT | s | d
	return p
}

func GetSender(p *Packet) string {

	s := uint32(p.Header) & MASK_SEND >> SHIFT_SEND
	return Packet_Section_name[int32(s)]
}
