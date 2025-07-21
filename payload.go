package wyoming

type Payloader interface {
	SetPayload(b []byte)
}

type Payload struct {
	payload []byte
}

func (p *Payload) SetPayload(b []byte) {
	p.payload = b
}

func (p *Payload) Payload() []byte {
	return p.payload
}
