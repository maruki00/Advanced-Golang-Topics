package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func initTracer(ctx context.Context) func() {

	client := otlptracegrpc.NewClient(otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint("localhost:4317"))
	exp, err := otlptrace.New(ctx, client)
	if err != nil {
		log.Fatalf("Failed to create OTLP exporter: %v", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("my-go-service"),
		)),
	)

	otel.SetTracerProvider(tp)

	return func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatalf("Failed to shutdown tracer provider: %v", err)
		}
	}
}

func main() {
	ctx := context.Background()
	cleanup := initTracer(ctx)
	defer cleanup()

	tracer := otel.Tracer("example-tracer")

	ctx, span := tracer.Start(ctx, "main-operation")
	defer span.End()

	processRequest(ctx)
}

func processRequest(ctx context.Context) {
	tracer := otel.Tracer("example-tracer")
	_, span := tracer.Start(ctx, "process-request")
	defer span.End()

	time.Sleep(500 * time.Second)
	fmt.Println("Processing request finished")
}
