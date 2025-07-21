package wyoming

const (
	TranscribeEvent      EventType = "transcribe"
	TranscriptEvent      EventType = "transcript"
	TranscriptStartEvent EventType = "transcript-start"
	TranscriptChunkEvent EventType = "transcript-chunk"
	TranscriptStopEvent  EventType = "transcript-stop"
)

type ASRProgram struct {
	Artifact

	Models    []ASRModel `json:"models"`
	Streaming bool       `json:"streaming"`
}

type ASRModel struct {
	Artifact

	Languages []string `json:"languages"`
}

type Transcribe struct {
	Name     string `json:"name,omitempty"`
	Language string `json:"language,omitempty"`
	Context  any    `json:"context,omitempty"`
}

func (t Transcribe) Event() Event {
	return Event{
		Type: TranscribeEvent,
		Data: t,
	}
}

type Transcript struct {
	Text     string `json:"text"`
	Language string `json:"language"`
	Context  any    `json:"context"`
}

type TranscriptStart struct {
	Language string `json:"language"`
	Context  any    `json:"context"`
}

type TranscriptChunk struct {
	Text string `json:"text"`
}

type TranscriptStop struct {
}
