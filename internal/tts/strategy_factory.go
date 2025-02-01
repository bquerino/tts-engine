package tts

import "errors"

type StrategyFactory struct {
	providers map[string]TTSProvider
}

func NewStrategyFactory() *StrategyFactory {
	return &StrategyFactory{
		providers: map[string]TTSProvider{
			"polly": &PollyTTS{},
		},
	}
}

func (sf *StrategyFactory) RegisterProvider(name string, provider TTSProvider) {
	sf.providers[name] = provider
}

func (sf *StrategyFactory) GetProvider(providerName string) (TTSProvider, error) {
	provider, exists := sf.providers[providerName]
	if !exists {
		return nil, errors.New("TTS Provider not found")
	}
	return provider, nil
}
