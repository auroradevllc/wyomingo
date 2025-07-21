package wyoming

// MicProgram is a microphone input
type MicProgram struct {
	Artifact

	Format AudioFormat `json:"mic_format"`
}
