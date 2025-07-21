package wyoming

const (
	RecognizeEvent     EventType = "recognize"
	IntentEvent        EventType = "intent"
	NotRecognizedEvent EventType = "not-recognized"
)

type IntentProgram struct {
	Artifact

	Models []IntentModel `json:"models"`
}

type IntentModel struct {
	Artifact

	Languages []string `json:"languages"`
}

type Entity struct {
}

type Recognize struct {
	Text    string  `json:"text"`
	Context Context `json:"context,omitempty"`
}

func (r Recognize) Event() Event {
	return Event{
		Type: RecognizeEvent,
		Data: r,
	}
}

type Intent struct {
	Name     string   `json:"name"`
	Entities []Entity `json:"entities"`
	Text     string   `json:"text,omitempty"`
	Context  Context  `json:"context,omitempty"`
}

func (i Intent) Event() Event {
	return Event{
		Type: IntentEvent,
	}
}

type NotRecognized struct {
	Text    string  `json:"text,omitempty"`
	Context Context `json:"context,omitempty"`
}
