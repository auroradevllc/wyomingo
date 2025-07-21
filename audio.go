package wyoming

import (
	"fmt"
	"github.com/gopxl/beep/v2"
)

type AudioStart struct {
	AudioFormat
	Timestamp int64 `json:"timestamp,omitempty"`
}

func (a AudioStart) Event() Event {
	return Event{
		Type: AudioStartEvent,
		Data: a,
	}
}

type AudioChunk struct {
	AudioFormat
	Audio     []byte `json:"-"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

func (c AudioChunk) SetPayload(b []byte) {
	c.Audio = b
}

func (c AudioChunk) Payload() []byte {
	return c.Audio
}

func (c AudioChunk) Event() Event {
	return Event{
		Type:    AudioChunkEvent,
		Data:    c,
		Payload: c.Audio,
	}
}

type AudioStop struct {
	Timestamp int64 `json:"timestamp,omitempty"`
}

func (s AudioStop) Event() Event {
	return Event{
		Type: AudioStopEvent,
		Data: s,
	}
}

type AudioFormat struct {
	Rate     int `json:"rate"`
	Width    int `json:"width"`
	Channels int `json:"channels"`
}

func StreamerToChunks(streamer beep.Streamer, format beep.Format) chan Eventable {
	events := make(chan Eventable)

	go func() {
		defer close(events)

		header := AudioFormat{
			Rate:     int(format.SampleRate),
			Channels: format.NumChannels,
			Width:    format.Width(),
		}

		events <- AudioStart{
			AudioFormat: header,
		}

		var (
			samples = make([][2]float64, 512)
			buffer  = make([]byte, len(samples)*format.Width())
		)

		for {
			n, ok := streamer.Stream(samples)
			if !ok {
				break
			}

			buf := buffer
			switch {
			case format.Precision == 1:
				for _, sample := range samples[:n] {
					buf = buf[format.EncodeUnsigned(buf, sample):]
				}
			case format.Precision == 2 || format.Precision == 3:
				for _, sample := range samples[:n] {
					buf = buf[format.EncodeSigned(buf, sample):]
				}
			default:
				panic(fmt.Errorf("wav: encode: invalid precision: %d", format.Precision))
			}

			events <- AudioChunk{
				AudioFormat: header,
				Audio:       buffer[:n*format.Width()],
			}
		}

		events <- AudioStop{}
	}()

	return events
}
