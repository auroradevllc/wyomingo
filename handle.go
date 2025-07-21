package wyoming

const (
	HandledEvent      EventType = "handled"
	NotHandledEvent   EventType = "not-handled"
	HandledStartEvent EventType = "handled-start"
	HandledChunkEvent EventType = "handled-chunk"
	HandledStopEvent  EventType = "handled-stop"
)

type HandleProgram struct {
	Artifact

	Models    []HandleModel `json:"models"`
	Streaming bool          `json:"supports_handled_streaming"`
}

type HandleModel struct {
	Artifact

	Languages []string `json:"languages"`
}

type Handled struct {
	Text    string  `json:"text,omitempty"`
	Context Context `json:"context,omitempty"`
}

type NotHandled struct {
	Text    string  `json:"text,omitempty"`
	Context Context `json:"context,omitempty"`
}

type HandledStart struct {
	Context Context `json:"context,omitempty"`
}

type HandledChunk struct {
	Text string `json:"text"`
}

type HandledStop struct {
}
