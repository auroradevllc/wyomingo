package wyoming

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestReceiveAudioEvents(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")

	if err != nil {
		t.Fatal(err)
	}

	defer ln.Close()

	go func() {
		conn, err := ln.Accept()

		if err != nil {
			return
		}

		defer conn.Close()

		data := `{"rate":22050,"width":2,"channels":1}`

		conn.Write([]byte(`{"type":"audio-start","data":{"rate":22050,"width":2,"channels":1}}` + "\n"))
		// Newer servers send data as a separate block after the header.
		conn.Write([]byte(`{"type":"audio-chunk","data_length":` + itoa(len(data)) + `,"payload_length":4}` + "\n"))
		conn.Write([]byte(data))
		conn.Write([]byte{1, 2, 3, 4})
		conn.Write([]byte(`{"type":"audio-stop"}` + "\n"))
	}()

	c, err := New("tcp://" + ln.Addr().String())

	if err != nil {
		t.Fatal(err)
	}

	defer c.Close()

	starts := make(chan *AudioStart, 1)
	chunks := make(chan *AudioChunk, 1)
	stops := make(chan *AudioStop, 1)

	c.AddSyncHandler(func(e *AudioStart) { starts <- e })
	c.AddSyncHandler(func(e *AudioChunk) { chunks <- e })
	c.AddSyncHandler(func(e *AudioStop) { stops <- e })

	go c.Run()

	timeout := time.After(5 * time.Second)

	select {
	case s := <-starts:
		if s.Rate != 22050 || s.Width != 2 || s.Channels != 1 {
			t.Fatalf("audio-start = %+v", s)
		}
	case <-timeout:
		t.Fatal("no audio-start")
	}

	select {
	case ch := <-chunks:
		if ch.Rate != 22050 || len(ch.Audio) != 4 || ch.Audio[3] != 4 {
			t.Fatalf("audio-chunk = %+v", ch)
		}
	case <-timeout:
		t.Fatal("no audio-chunk")
	}

	select {
	case <-stops:
	case <-timeout:
		t.Fatal("no audio-stop")
	}
}

func TestNewWithContextHonorsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := NewWithContext(ctx, "tcp://127.0.0.1:1"); err == nil {
		t.Fatal("expected an error from a cancelled context")
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}

	var b []byte

	for ; n > 0; n /= 10 {
		b = append([]byte{byte('0' + n%10)}, b...)
	}

	return string(b)
}
