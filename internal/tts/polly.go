package tts

import (
	"context"
	"fmt"
	"os"
	"tts-engine/internal/config"
	"tts-engine/internal/monitoring"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/aws/aws-sdk-go-v2/service/polly/types"
)

type PollyTTS struct{}

func (p *PollyTTS) GenerateSpeech(text, language, voice string) (string, error) {
	// Load AWS configurations
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithRegion(config.AppConfig.AWSRegion),
	)
	if err != nil {
		monitoring.ErrorLog("Failed to load AWS configuration", err, nil)
		return "", err
	}

	client := polly.NewFromConfig(cfg)

	input := &polly.SynthesizeSpeechInput{
		Text:         &text,
		OutputFormat: types.OutputFormatMp3,
		VoiceId:      types.VoiceId(voice),
	}

	resp, err := client.SynthesizeSpeech(context.TODO(), input)
	if err != nil {
		monitoring.ErrorLog("Failed to call Polly", err, map[string]interface{}{
			"text":     text,
			"language": language,
			"voice":    voice,
		})
		return "", err
	}

	// Save the audio to a temporary file TODO: Send to S3
	filePath := fmt.Sprintf("/tmp/%s.mp3", voice)
	file, err := os.Create(filePath)
	if err != nil {
		monitoring.ErrorLog("Failed to create audio file", err, map[string]interface{}{
			"file_path": filePath,
		})
		return "", err
	}
	defer file.Close()

	_, err = file.ReadFrom(resp.AudioStream)
	if err != nil {
		monitoring.ErrorLog("Failed to save audio", err, map[string]interface{}{
			"file_path": filePath,
		})
		return "", err
	}

	monitoring.InfoLog("Audio file successfully generated", map[string]interface{}{
		"file_path": filePath,
		"voice":     voice,
	})

	return filePath, nil
}
