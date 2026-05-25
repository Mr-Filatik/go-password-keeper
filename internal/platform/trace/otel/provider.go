package oteltrace

import (
	"context"
	"fmt"

	"github.com/mr-filatik/go-password-keeper/internal/platform/log"
	mytrace "github.com/mr-filatik/go-password-keeper/internal/platform/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.5.0"
)

type TracerProviderConfig struct {
	// ExporterType enum
	InsecureMode bool   // включен только на dev
	Endpoint     string // "localhost:4317"
	Metadata     TracerProviderConfigMetadata
}

type TracerProviderConfigMetadata struct {
	ServiceName string
	Version     string
	Environment string
}

// Tracer(name string, options ...TracerOption) Tracer

func NewTracerProvider(ctx context.Context, cfg TracerProviderConfig, logger log.ILogger) (*TracerProvider, error) {
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(), // Отключаем TLS для локальной разработки
		otlptracegrpc.WithEndpoint(cfg.Endpoint),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP gRPC exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.Metadata.ServiceName),
			semconv.ServiceVersionKey.String(cfg.Metadata.Version),
			semconv.DeploymentEnvironmentKey.String(cfg.Metadata.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // Логируем абсолютно все трассы
		sdktrace.WithBatcher(exporter),                // Асинхронная буферизированная отправка спанов
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetLogger(NewOTelLoggerAdapter(logger))
	otel.SetErrorHandler(NewOTelErrorHandler(logger))

	return &TracerProvider{
		TracerProvider: tp,
	}, nil
}

var otp mytrace.ITracerProvider = &TracerProvider{}

type TracerProvider struct {
	*sdktrace.TracerProvider
}

func (p *TracerProvider) Tracer(name string) mytrace.ITracer {
	tracer := p.TracerProvider.Tracer(name)

	return &Tracer{
		Tracer: tracer,
	}
}
