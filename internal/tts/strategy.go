package tts

type TTSProvider interface {
	GenerateSpeech(text, language, voice string) (string, error)
}
