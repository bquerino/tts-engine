package tts

import (
	"tts-engine/internal/monitoring"
	"tts-engine/internal/repository"
)

// Balancer is responsible for selecting the best TTS provider
type Balancer struct {
	repo repository.UsageRepository
}

// NewBalancer initializes the load balancer with a usage repository
func NewBalancer(repo repository.UsageRepository) *Balancer {
	monitoring.InfoLog("Initializing TTS load balancer...", nil)
	return &Balancer{repo: repo}
}

// SelectProvider chooses the best provider based on current usage
func (b *Balancer) SelectProvider() string {
	providers := []string{"polly", "google", "azure"}

	for _, provider := range providers {
		messages, characters, err := b.repo.GetUsage(provider)
		if err == nil && messages < 10000 && characters < 5000000 {
			monitoring.InfoLog("Selected TTS provider", map[string]interface{}{
				"provider":   provider,
				"messages":   messages,
				"characters": characters,
			})
			return provider
		}
	}

	monitoring.WarnLog("⚠️ All providers have reached their limits! Using Polly as a fallback.", nil)
	return "polly"
}

// TrackUsage updates the usage statistics for a provider
func (b *Balancer) TrackUsage(provider string, text string) error {
	err := b.repo.IncrementMessageCount(provider)
	if err != nil {
		monitoring.ErrorLog("Failed to update message count", err, map[string]interface{}{
			"provider": provider,
		})
		return err
	}

	err = b.repo.IncrementCharacterCount(provider, len(text))
	if err != nil {
		monitoring.ErrorLog("Failed to update character count", err, map[string]interface{}{
			"provider": provider,
			"length":   len(text),
		})
		return err
	}

	monitoring.InfoLog("Provider usage updated", map[string]interface{}{
		"provider":   provider,
		"characters": len(text),
	})

	return nil
}
