package config

import (
	"os"
	"tts-engine/internal/monitoring"

	"github.com/joho/godotenv"
)

// Config representa as configurações da aplicação
type Config struct {
	Port      string
	AWSRegion string
	AWSKey    string
	AWSSecret string
}

var AppConfig Config

func Load() {
	if err := godotenv.Load(); err != nil {
		monitoring.WarnLog("File .env not found, loading environment variables.", nil)
	}

	AppConfig = Config{
		Port:      getEnv("PORT", "8080"),
		AWSRegion: getEnv("AWS_REGION", "us-east-1"),
		AWSKey:    getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecret: getEnv("AWS_SECRET_ACCESS_KEY", ""),
	}

	monitoring.InfoLog("Success on load configs!", map[string]interface{}{
		"port":       AppConfig.Port,
		"aws_region": AppConfig.AWSRegion,
	})
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
