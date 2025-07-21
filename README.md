Wyoming for Go
=

An low-level implementation of the [Wyoming](https://github.com/OHF-Voice/wyoming) protocol for Go.

This implementation is "incomplete" (in that it doesn't support all events going out/in), though
it's easy to add and going to be done slowly.

How it works
------------

All events are defined as structs. There is currently no way to write to the extra data body defined by the protocol,
however decoding this is supported.

For outgoing events, they inherit `Eventable`:

```go
type Eventable interface {
	Event() Event
}
```

This allows them to be serialized back with a `Type`, `Data`, and optional `Payload` attribute.

For incoming events, they are defined in `eventTypes` under `wyoming.go`, which allows them
to be created dynamically for incoming events. After they are decoded, they are passed to the
[handler](https://github.com/auroradevllc/handler) package, which is a reflection-based event handler.

Audio Handling
--------------

This library uses [beep](https://github.com/gopxl/beep) for audio handling. You do not have to use beep, any method of
generating PCM audio will work. The library does no processing of return audio.

Any convenience functionality will be written around beep, such as `StreamerToChunks`.

Example
-------

An example/test of `whisper.cpp` is defined in `asr_test.go`, though it does not contain the necessary
wav files. This also shows how to write samples out to Audio Chunks, which may be implemented as a helper
in the future via channels.

```go
package main

import (
	"context"
	"log"
	"github.com/auroradevllc/wyomingo"
)

func main() {
	// Create client
	c, err := wyoming.New("tcp://127.0.0.1:10300")

	if err != nil {
		log.Fatal("Unable to connect to wyoming server:", err)
	}

	// Start processing messages
	go c.Run()

	// Write Describe event, this is required.
	err = c.Write(wyoming.Describe{})

	if err != nil {
		log.Fatal("Unable to write describe event:", err)
	}

	// Wait for the Info event. Note that events are pointers.
	// See the 'handler" package in the same organization for handler-specific functions
	info := c.WaitFor(context.Background(), func(e any) bool {
		_, ok := e.(*Info)
		return ok
	})

	log.Println("Info event received: ", info)
}
```