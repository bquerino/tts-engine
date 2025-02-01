package tts

import (
	"fmt"
	"sync"
	"tts-engine/internal/monitoring"
)

type ProviderUsage struct {
	Messages   int
	Characters int
}

type Quotas struct {
	MaxMessages   int
	MaxCharacters int
}

var (
	usage  sync.Map
	quotas sync.Map
)

// InitQuotas initializes the quotas and usage tracking for providers
// TODO: Move this logic to a repository
func InitQuotas() {
	quotas.Store("polly", &Quotas{MaxMessages: 10000, MaxCharacters: 5000000})
	quotas.Store("google", &Quotas{MaxMessages: 10000, MaxCharacters: 5000000})
	quotas.Store("azure", &Quotas{MaxMessages: 10000, MaxCharacters: 5000000})

	usage.Store("polly", &ProviderUsage{})
	usage.Store("google", &ProviderUsage{})
	usage.Store("azure", &ProviderUsage{})

	monitoring.InfoLog("Quotas initialized successfully", nil)
}

// TrackUsage registers the usage of a provider and checks quota limits
func TrackUsage(provider string, text string) error {
	// Retrieve provider usage
	uValue, exists := usage.Load(provider)
	if !exists {
		err := fmt.Errorf("unknown provider: %s", provider)
		monitoring.ErrorLog("Failed to track usage", err, map[string]interface{}{
			"provider": provider,
		})
		return err
	}
	u := uValue.(*ProviderUsage)

	// Retrieve provider quota
	qValue, exists := quotas.Load(provider)
	if !exists {
		err := fmt.Errorf("quota not found for provider: %s", provider)
		monitoring.ErrorLog("Quota not found", err, map[string]interface{}{
			"provider": provider,
		})
		return err
	}
	q := qValue.(*Quotas)

	// Update usage counts
	u.Messages++
	u.Characters += len(text)

	// Check quota limits
	if u.Messages > q.MaxMessages || u.Characters > q.MaxCharacters {
		err := fmt.Errorf("usage limit exceeded for provider: %s", provider)
		monitoring.WarnLog("Provider usage limit exceeded", map[string]interface{}{
			"provider":   provider,
			"messages":   u.Messages,
			"characters": u.Characters,
		})
		return err
	}

	monitoring.InfoLog("Provider usage updated", map[string]interface{}{
		"provider":   provider,
		"messages":   u.Messages,
		"characters": u.Characters,
	})

	return nil
}
