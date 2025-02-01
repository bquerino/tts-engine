package api

import (
	"tts-engine/internal/monitoring"
	"tts-engine/internal/tts"

	"github.com/gofiber/fiber/v2"
)

func SynthesizeHandler(balancer *tts.Balancer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req SynthesizeRequest
		if err := c.BodyParser(&req); err != nil {
			monitoring.ErrorLog("Error to process this request.", err, nil)
			return c.Status(fiber.StatusBadRequest).JSON(SynthesizeResponse{
				Error: "Error when processing the request",
			})
		}

		monitoring.InfoLog("Received request", map[string]interface{}{
			"method": c.Method(),
			"path":   c.Path(),
			"text":   req.Text,
		})

		if req.Provider == "" {
			req.Provider = balancer.SelectProvider()
		}

		err := balancer.TrackUsage(req.Provider, req.Text)
		if err != nil {
			return c.Status(fiber.StatusTooManyRequests).JSON(SynthesizeResponse{
				Error: "Voice synhtesis limit reached",
			})
		}

		factory := tts.NewStrategyFactory()
		ttsProvider, err := factory.GetProvider(req.Provider)
		if err != nil {
			monitoring.ErrorLog("Error to select synth provider", err, map[string]interface{}{
				"provider": req.Provider,
			})
			return c.Status(fiber.StatusBadRequest).JSON(SynthesizeResponse{
				Error: "Invalid TTS provider",
			})
		}

		audioURL, err := ttsProvider.GenerateSpeech(req.Text, req.Language, req.Voice)
		if err != nil {
			monitoring.ErrorLog("Error on audio creation", err, map[string]interface{}{
				"provider": req.Provider,
			})
			return c.Status(fiber.StatusInternalServerError).JSON(SynthesizeResponse{
				Error: "Error on audio creation",
			})
		}

		monitoring.InfoLog("Success on audio creation", map[string]interface{}{
			"provider": req.Provider,
			"audioURL": audioURL,
		})

		return c.JSON(SynthesizeResponse{
			AudioURL: audioURL,
		})
	}
}
