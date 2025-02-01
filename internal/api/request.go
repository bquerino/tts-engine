package api

type SynthesizeRequest struct {
	Text     string `json:"text" validate:"required"`
	Language string `json:"language" validate:"required"`
	Voice    string `json:"voice" validate:"required"`
	Provider string `json:"provider,omitempty"`
}
