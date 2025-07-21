package wyoming

type Describe struct {
}

func (d Describe) Event() Event {
	return Event{
		Type: DescribeEvent,
	}
}

type Info struct {
	ASR    []ASRProgram    `json:"asr"`
	TTS    []TTSProgram    `json:"tts"`
	Handle []HandleProgram `json:"handle"`
	Intent []IntentProgram `json:"intent"`
	Wake   []WakeProgram   `json:"wake"`
	Mic    []MicProgram    `json:"mic"`
	Snd    []SndProgram    `json:"snd"`
}

type Artifact struct {
	Name        string      `json:"name"`
	Attribution Attribution `json:"attribution"`
	Installed   bool        `json:"installed"`
	Description string      `json:"description"`
	Version     string      `json:"version"`
}
