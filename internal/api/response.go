package api

type SynthesizeResponse struct {
	AudioURL string `json:"audio_url"`
	Error    string `json:"error,omitempty"`
}
