package main

import (
	"log"
	"tts-engine/internal/api"
	"tts-engine/internal/monitoring"
	"tts-engine/internal/repository"
	"tts-engine/internal/storage"
	"tts-engine/internal/tts"

	"github.com/gofiber/fiber/v2"
)

func main() {
	monitoring.InitLogger()

	_, err := monitoring.InitTracer()
	if err != nil {
		log.Fatalf("Error on init tracing: %v", err)
	}

	storage.InitRedis()

	repo := repository.NewRedisUsageRepository(storage.RedisClient)

	balancer := tts.NewBalancer(repo)

	app := fiber.New()

	app.Use(monitoring.TracingMiddleware())

	app.Post("/synthesize", api.SynthesizeHandler(balancer))

	monitoring.InfoLog("🚀 Server running on port 8080...", nil)
	log.Fatal(app.Listen(":8080"))
}
