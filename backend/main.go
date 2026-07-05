package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RicardoSandoval11/apartamentos/backend/constants"
	"github.com/RicardoSandoval11/apartamentos/backend/middleware"
	"github.com/RicardoSandoval11/apartamentos/backend/pkg/apartment"

	httptransport "github.com/go-kit/kit/transport/http"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func main() {

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer cancel()

	// 2. Configuraciones de logging
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	).With("service", constants.SERVICE_NAME)

	slog.SetDefault(logger)

	// 3. Configuracion de OpenTelemetry
	exporter, _ := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
	metricExporter, _ := otlpmetricgrpc.New(ctx)
	res, _ := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceName(constants.SERVICE_NAME),
			semconv.ServiceVersion(constants.SERVICE_VERSION),
			attribute.String("environment", constants.ENVIRONMENT),
		),
	)

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	metricProvider := metric.NewMeterProvider(
		metric.WithReader(
			metric.NewPeriodicReader(metricExporter),
		),
		metric.WithResource(res),
	)

	otel.SetTracerProvider(traceProvider)
	otel.SetMeterProvider(metricProvider)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	// GO KIT INTEGRATION
	aptService := apartment.NewApartmentService()

	aptEndpoint := apartment.MakeGetApartmentsEndpoint(aptService)
	{
		aptEndpoint = middleware.LoggingMiddleware()(aptEndpoint)
		aptEndpoint = middleware.AuthMiddleware()(aptEndpoint)
	}

	aptHandler := httptransport.NewServer(
		aptEndpoint,
		apartment.DecodeGetApartmentRequest,
		apartment.EncodeGetApartmentResponse,
		httptransport.ServerBefore(middleware.ExtractTokenFromHeader),
	)

	mux := http.NewServeMux()
	mux.Handle("/apartment", aptHandler)

	handler := otelhttp.NewHandler(mux, "api-server")

	errs := make(chan error)

	srv := &http.Server{
		Addr:    constants.APPLICATION_PORT,
		Handler: handler,
	}

	go func() {
		errs <- srv.ListenAndServe()
	}()

	select {
	case e := <-errs:
		{
			slog.Error(
				"Application was not able to start",
				"error", e,
			)
		}
	case <-ctx.Done():
		{
			slog.Info("Shutting down")

			shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			_ = traceProvider.Shutdown(shutdownCtx)
			_ = metricProvider.Shutdown(shutdownCtx)
			if err := srv.Shutdown(shutdownCtx); err != nil {
				slog.Error("Shutdown failed", "error", err)
			}
		}
	}
}
