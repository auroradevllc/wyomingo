package wyoming

import (
	"encoding/json"
)

type EventType string

const (
	DescribeEvent   EventType = "describe"
	InfoEvent       EventType = "info"
	AudioStartEvent EventType = "audio-start"
	AudioChunkEvent EventType = "audio-chunk"
	AudioStopEvent  EventType = "audio-stop"
)

type EventHeader struct {
	Type          EventType       `json:"type"`
	DataLength    int             `json:"data_length,omitempty"`
	PayloadLength int             `json:"payload_length,omitempty"`
	Data          json.RawMessage `json:"data,omitempty"`
}

type IncomingEvent struct {
	Type          EventType       `json:"type"`
	DataLength    int             `json:"data_length,omitempty"`
	PayloadLength int             `json:"payload_length,omitempty"`
	Data          json.RawMessage `json:"data,omitempty"`
}

type Event struct {
	Type          EventType `json:"type"`
	DataLength    int       `json:"data_length,omitempty"`
	PayloadLength int       `json:"payload_length,omitempty"`
	Data          any       `json:"data"`
	Payload       []byte    `json:"-"`
}

type Eventable interface {
	Event() Event
}
