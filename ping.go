package wyoming

const (
	PingEvent EventType = "ping"
	PongEvent EventType = "pong"
)

type Ping struct {
	Text string `json:"text,omitempty"`
}

func (p Ping) Event() Event {
	return Event{
		Type: PingEvent,
		Data: p,
	}
}

type Pong struct {
	Text string `json:"text,omitempty"`
}

func (p Pong) Event() Event {
	return Event{
		Type: PongEvent,
		Data: p,
	}
}
