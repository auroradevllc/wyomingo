package wyoming

const (
	SynthesizeEvent        EventType = "synthesize"
	SynthesizeStartEvent   EventType = "synthesize-start"
	SynthesizeChunkEvent   EventType = "synthesize-chunk"
	StopSynthesizeEvent    EventType = "synthesize-stop"
	SynthesizeStoppedEvent EventType = "synthesize-stopped"
)

type TTSProgram struct {
	Artifact
	Streaming bool `json:"supports_synthesize_streaming"`
}

type TTSVoice struct {
	Artifact

	Languages []string          `json:"languages"`
	Speakers  []TTSVoiceSpeaker `json:"speakers"`
}

type TTSVoiceSpeaker struct {
	Name string `json:"name"`
}

type Synthesize struct {
	Text    string    `json:"text"`
	Voice   *TTSVoice `json:"voice,omitempty"`
	Context Context   `json:"context,omitempty"`
}

func (s Synthesize) Event() Event {
	return Event{
		Type: SynthesizeEvent,
		Data: s,
	}
}

type SynthesizeStart struct {
	Voice   *TTSVoice `json:"voice,omitempty"`
	Context Context   `json:"context,omitempty"`
}

func (s SynthesizeStart) Event() Event {
	return Event{
		Type: SynthesizeStartEvent,
		Data: s,
	}
}

type SynthesizeChunk struct {
	Text string `json:"text"`
}

func (s SynthesizeChunk) Event() Event {
	return Event{
		Type: SynthesizeChunkEvent,
		Data: s,
	}
}

type SynthesizeStop struct {
}

func (s SynthesizeStop) Event() Event {
	return Event{
		Type: StopSynthesizeEvent,
	}
}

type SynthesizeStopped struct {
}

func (s SynthesizeStopped) Event() Event {
	return Event{
		Type: SynthesizeStoppedEvent,
	}
}
