package wyoming

// SndProgram is a sound output program
type SndProgram struct {
	Artifact

	Format AudioFormat `json:"snd_format"`
}
