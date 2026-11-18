package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RicardoSandoval11/apartamentos/identity-service/conf"
	"github.com/RicardoSandoval11/apartamentos/identity-service/db"
	"github.com/RicardoSandoval11/apartamentos/identity-service/db/migrations"
	"github.com/RicardoSandoval11/apartamentos/identity-service/middleware"
	"github.com/RicardoSandoval11/apartamentos/identity-service/pkg/apartment"
	"github.com/RicardoSandoval11/apartamentos/identity-service/pkg/identity"
	"github.com/RicardoSandoval11/apartamentos/identity-service/pkg/verification"

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

	cfg, err := conf.GetEnv()

	if err != nil {
		slog.Error("could not load environment variables", "error", err)
		os.Exit(1)
	}

	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	).With("service", cfg.ServiceName)

	slog.SetDefault(logger)

	exporter, _ := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
	metricExporter, _ := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithInsecure())
	res, _ := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			attribute.String("environment", cfg.Environment),
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

	dbInstance := db.GetDatabase(cfg.DbConnectionString)

	logger.Info("checking for pending migrations...")

	if err := migrations.Run(dbInstance); err != nil {
		logger.Error("failed applying mgirations", "error", err)
		os.Exit(1)
	}

	identityRepository := identity.NewIdentityRepository(dbInstance)
	verificationRepository := verification.NewVerificationRepository(dbInstance)

	identityService := identity.NewIdentityService(identityRepository, verificationRepository)

	identityEndpoint := apartment.MakeGetApartmentsEndpoint(identityService)
	{
		identityEndpoint = middleware.LoggingMiddleware()(identityEndpoint)
		identityEndpoint = middleware.AuthMiddleware()(identityEndpoint)
	}

	identityHandler := httptransport.NewServer(
		identityEndpoint,
		identity.DecodeGetApartmentRequest,
		identity.EncodeGetApartmentResponse,
		httptransport.ServerBefore(middleware.ExtractTokenFromHeader),
	)

	mux := http.NewServeMux()
	mux.Handle("/identity", identityHandler)

	handler := otelhttp.NewHandler(mux, "api-server")

	errs := make(chan error)

	srv := &http.Server{
		Addr:    cfg.ApplicationPort,
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
