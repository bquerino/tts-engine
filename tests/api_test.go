package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"tts-engine/internal/api"
	"tts-engine/internal/tts"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTTSProvider struct {
	mock.Mock
}

// GenerateSpeech mocks the speech synthesis function
func (m *MockTTSProvider) GenerateSpeech(text, language, voice string) (string, error) {
	args := m.Called(text, language, voice)
	return args.String(0), args.Error(1)
}

// TestSynthesizeHandler tests the /synthesize endpoint
func TestSynthesizeHandler(t *testing.T) {
	app := fiber.New()

	// Mock provider
	mockTTS := new(MockTTSProvider)
	mockTTS.On("GenerateSpeech", "Hello", "en-US", "Matthew").Return("https://example.com/audio.mp3", nil)

	// Create and register the provider factory
	ttsFactory := tts.NewStrategyFactory()
	ttsFactory.RegisterProvider("polly", mockTTS)

	// Define the /synthesize endpoint
	app.Post("/synthesize", func(c *fiber.Ctx) error {
		var req api.SynthesizeRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(http.StatusBadRequest).JSON(api.SynthesizeResponse{
				Error: "Failed to process request",
			})
		}

		ttsProvider, err := ttsFactory.GetProvider(req.Provider)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(api.SynthesizeResponse{
				Error: "Invalid TTS provider",
			})
		}

		audioURL, err := ttsProvider.GenerateSpeech(req.Text, req.Language, req.Voice)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(api.SynthesizeResponse{
				Error: "Error generating audio",
			})
		}

		return c.JSON(api.SynthesizeResponse{
			AudioURL: audioURL,
		})
	})

	// Create request body
	reqBody, _ := json.Marshal(api.SynthesizeRequest{
		Text:     "Hello",
		Language: "en-US",
		Voice:    "Matthew",
		Provider: "polly",
	})
	req := httptest.NewRequest(http.MethodPost, "/synthesize", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, _ := app.Test(req, -1)

	// Validate response status
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Validate JSON response
	var jsonResponse api.SynthesizeResponse
	_ = json.NewDecoder(resp.Body).Decode(&jsonResponse)
	assert.Equal(t, "https://example.com/audio.mp3", jsonResponse.AudioURL)

	// Ensure mock expectations were met
	mockTTS.AssertExpectations(t)
}
