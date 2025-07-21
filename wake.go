package wyoming

const (
	DetectEvent      EventType = "detect"
	DetectionEvent   EventType = "detection"
	NotDetectedEvent EventType = "not-detected"
)

// WakeProgram is a wake word detection service
type WakeProgram struct {
	Artifact

	Models []WakeModel `json:"models"`
}

// WakeModel is a wake word detection model
type WakeModel struct {
	Artifact

	Languages []string `json:"languages"`
	Phrase    string   `json:"phrase"`
}

type Detect struct {
	Names   []string `json:"names,omitempty"`
	Context Context  `json:"context,omitempty"`
}

func (d Detect) Event() Event {
	return Event{
		Type: DetectEvent,
		Data: d,
	}
}

type Detection struct {
	Name      string  `json:"name"`
	Timestamp int64   `json:"timestamp"`
	Speaker   string  `json:"speaker"`
	Context   Context `json:"context,omitempty"`
}

type NotDetected struct {
	Context Context `json:"context,omitempty"`
}
