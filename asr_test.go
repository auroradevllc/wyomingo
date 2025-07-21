package wyoming

import (
	"context"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/wav"
	"os"
	"testing"
)

func TestASR(t *testing.T) {
	c, err := New("tcp://127.0.0.1:10300")

	if err != nil {
		t.Fatal(err)
	}

	go c.Run()

	err = c.Write(Describe{})

	if err != nil {
		t.Fatal("Unable to write:", err)
	}

	info := c.WaitFor(context.Background(), func(e any) bool {
		_, ok := e.(*Info)
		return ok
	})

	t.Log("Info event received: ", info)

	ch, done := c.ChanFor(func(ev any) bool {
		return true
	})

	defer done()

	err = c.Write(Transcribe{})

	if err != nil {
		t.Fatal("Unable to write:", err)
	}

	f, err := os.Open("tests_turn_on_the_living_room_lamp.wav")

	if err != nil {
		t.Fatal(err)
	}

	streamer, format, err := wav.Decode(f)

	if err != nil {
		t.Fatal(err)
	}

	newSampleRate := beep.SampleRate(16000)

	resampler := beep.Resample(4, format.SampleRate, newSampleRate, streamer)

	newFormat := beep.Format{
		SampleRate:  newSampleRate,
		NumChannels: 1,
		Precision:   format.Precision,
	}

	chunks := StreamerToChunks(resampler, newFormat)

	_ = c.WriteChan(chunks)

	t.Log("Waiting for event")

	for {
		e, ok := <-ch

		if !ok {
			t.Fatal("Unexpected event")
		}

		if transcript, ok := e.(*Transcript); ok {
			t.Log("Got transcript", transcript.Text)
			break
		}
	}
}
