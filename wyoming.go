package wyoming

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/auroradevllc/handler"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"net/url"
	"sync"
)

type EventConstructor func() any

func eventStructFunc[V any]() EventConstructor {
	return func() any {
		return new(V)
	}
}

var (
	ErrUnsupportedScheme = errors.New("unsupported connection type")

	eventTypes = map[EventType]EventConstructor{
		// Base events
		InfoEvent: eventStructFunc[Info](),

		// ASR Events
		TranscriptEvent:      eventStructFunc[Transcript](),
		TranscriptStartEvent: eventStructFunc[TranscriptStart](),
		TranscriptChunkEvent: eventStructFunc[TranscriptChunk](),
		TranscriptStopEvent:  eventStructFunc[TranscriptStop](),

		// TTS Events
		SynthesizeStoppedEvent: eventStructFunc[SynthesizeStopped](),

		// Ping Events
		PingEvent: eventStructFunc[Ping](),
		PongEvent: eventStructFunc[Pong](),

		// Wake Word Events
		DetectionEvent:   eventStructFunc[Detection](),
		NotDetectedEvent: eventStructFunc[NotDetected](),

		// Intent Events
		IntentEvent:        eventStructFunc[Intent](),
		NotRecognizedEvent: eventStructFunc[NotRecognized](),

		// Handle Events
		HandledEvent:      eventStructFunc[Handled](),
		NotHandledEvent:   eventStructFunc[NotHandled](),
		HandledStartEvent: eventStructFunc[HandledStart](),
		HandledChunkEvent: eventStructFunc[HandledChunk](),
		HandledStopEvent:  eventStructFunc[HandledStop](),

		// Timer Events
		TimerStartedEvent:   eventStructFunc[TimerStarted](),
		TimerUpdatedEvent:   eventStructFunc[TimerUpdated](),
		TimerCancelledEvent: eventStructFunc[TimerCancelled](),
		TimerFinishedEvent:  eventStructFunc[TimerFinished](),
	}
)

type Context map[string]any

type Client interface {
	handler.HandlerInterface
	Write(event Eventable) error
	WriteChan(events chan Eventable) error
	WriteMultiple(events ...Eventable) error
	Run()
}

type client struct {
	*handler.Handler
	uri       string
	reader    *bufio.Reader
	writer    *bufio.Writer
	closer    io.Closer
	writeLock *sync.Mutex
}

func New(uri string) (Client, error) {
	c := &client{
		Handler:   handler.New(),
		uri:       uri,
		writeLock: &sync.Mutex{},
	}

	if err := c.Connect(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *client) Connect() error {
	u, err := url.Parse(c.uri)

	if err != nil {
		return err
	}

	switch u.Scheme {
	case "tcp":
		conn, err := net.Dial("tcp", u.Host)

		if err != nil {
			return err
		}

		c.reader = bufio.NewReader(conn)
		c.writer = bufio.NewWriter(conn)
		c.closer = conn
	case "unix":
		conn, err := net.Dial("unix", u.Path)

		if err != nil {
			return err
		}

		c.reader = bufio.NewReader(conn)
		c.writer = bufio.NewWriter(conn)
		c.closer = conn
	default:
		return ErrUnsupportedScheme
	}

	return nil
}

func (c *client) Run() {
	for {
		line, err := c.reader.ReadString('\n')

		if err != nil {
			return
		}

		var header IncomingEvent

		if err := json.Unmarshal([]byte(line), &header); err != nil {
			return
		}

		var data, payload []byte

		if header.DataLength > 0 {
			data = make([]byte, header.DataLength)

			if _, err := io.ReadFull(c.reader, data); err != nil {
				// Short read
				return
			}
		}

		if header.PayloadLength > 0 {
			payload = make([]byte, header.PayloadLength)

			if _, err := io.ReadFull(c.reader, payload); err != nil {
				return
			}
		}

		c.handleEvent(header, data, payload)
	}
}

func (c *client) handleEvent(header IncomingEvent, data, payload []byte) {
	log.WithFields(log.Fields{
		"type": header.Type,
	}).Debug("Handle event")
	var v any

	if c, ok := eventTypes[header.Type]; ok {
		v = c()
	}

	if v == nil {
		log.WithField("type", header.Type).Warn("No handler for event")
		return
	}

	// Decode original data
	if header.Data != nil {
		if err := json.Unmarshal(header.Data, v); err != nil {
			log.WithError(err).Error("Failed to unmarshal event data")
			return
		}
	}

	// Merge data into the struct
	if data != nil {
		if err := json.Unmarshal(data, v); err != nil {
			log.WithError(err).Error("Failed to unmarshal payload data")
			return
		}
	}

	if payload != nil {
		if p, ok := v.(Payloader); ok {
			p.SetPayload(payload)
		}
	}

	c.Call(v)
}

func (c *client) Write(e Eventable) error {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()

	event := e.Event()

	// Populate payload length
	if event.Payload != nil {
		event.PayloadLength = len(event.Payload)
	}

	b, err := json.Marshal(event)

	if err != nil {
		return err
	}

	_, err = c.writer.Write(b)

	if err != nil {
		return err
	}

	_, err = c.writer.WriteRune('\n')

	if err != nil {
		return err
	}

	if event.Payload != nil {
		_, err = c.writer.Write(event.Payload)

		if err != nil {
			return err
		}
	}

	return c.writer.Flush()
}

func (c *client) WriteChan(events chan Eventable) error {
	for e := range events {
		if err := c.Write(e); err != nil {
			return err
		}
	}

	return nil
}

func (c *client) WriteMultiple(events ...Eventable) error {
	for _, event := range events {
		if err := c.Write(event); err != nil {
			return err
		}
	}

	return nil
}

func (c *client) Close() error {
	return c.closer.Close()
}
