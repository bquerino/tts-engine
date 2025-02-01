package monitoring

import (
	"log/slog"
	"os"
)

// TODO: create library to abstract logger
var Logger *slog.Logger

func InitLogger() {
	// default: "json", options: "json" or "text")
	logFormat := os.Getenv("LOG_FORMAT")
	if logFormat == "" {
		logFormat = "json"
	}

	var handler slog.Handler
	switch logFormat {
	case "text":
		handler = slog.NewTextHandler(os.Stdout, nil)
	default:
		handler = slog.NewJSONHandler(os.Stdout, nil)
	}

	Logger = slog.New(handler)
}

func InfoLog(message string, fields map[string]interface{}) {
	if Logger == nil {
		InitLogger()
	}
	Logger.Info(message, convertFields(fields)...)
}

func WarnLog(message string, fields map[string]interface{}) {
	if Logger == nil {
		InitLogger()
	}
	Logger.Warn(message, convertFields(fields)...)
}

func ErrorLog(message string, err error, fields map[string]interface{}) {
	if Logger == nil {
		InitLogger()
	}
	Logger.Error(message, append(convertFields(fields), "error", err)...)
}

func convertFields(fields map[string]interface{}) []any {
	var result []any
	for k, v := range fields {
		result = append(result, k, v)
	}
	return result
}
