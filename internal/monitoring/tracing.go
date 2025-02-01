package monitoring

import (
	"context"
	"os"

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

// InitTracer configures the OpenTelemetry Tracer Provider
func InitTracer() (*trace.TracerProvider, error) {
	otelEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if otelEndpoint == "" {
		otelEndpoint = "http://otel-collector:4317" // Default for Jaeger or Tempo
	}

	exporter, err := otlptrace.New(context.Background(), otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(otelEndpoint),
		otlptracehttp.WithInsecure(),
	))
	if err != nil {
		ErrorLog("Failed to configure OLTP exporter", err, map[string]interface{}{
			"otel_endpoint": otelEndpoint,
		})
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("tts-engine"),
		)),
	)

	otel.SetTracerProvider(tp)

	InfoLog("Tracing successfully configured!", map[string]interface{}{
		"otel_endpoint": otelEndpoint,
	})
	return tp, nil
}

// TracingMiddleware adds the tracing middleware to Fiber
func TracingMiddleware() fiber.Handler {
	return otelfiber.Middleware()
}
