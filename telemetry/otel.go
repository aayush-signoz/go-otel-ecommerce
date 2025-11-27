package telemetry

import (
	"context"
	"os"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/bridges/otellogrus"
	otelruntime "go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	otel_log "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"

	"go-otel-ecommerce/config"
)

var (
	OrdersTotal         metric.Int64Counter
	ProductOrderCounter metric.Int64Counter
	HttpRequestCount    metric.Int64Counter
	HttpDurationBucket  metric.Float64Histogram
	GoroutinesGauge     metric.Int64ObservableGauge
	MemoryGauge         metric.Int64ObservableGauge
)

func otelResource() *resource.Resource {
	hostname, _ := os.Hostname()
	res, _ := resource.New(context.Background(),
		resource.WithSchemaURL("https://opentelemetry.io/schemas/1.11.0"),
		resource.WithAttributes(
			attribute.String("service.name", config.ServiceName),
			attribute.String("service.version", "1.0.0"),
			attribute.String("service.instance.id", hostname),
		),
	)
	return res
}

func InitTracer() func(context.Context) error {
	var opt otlptracegrpc.Option
	if config.Insecure == "false" {
		opt = otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	} else {
		opt = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptracegrpc.New(context.Background(),
		opt,
		otlptracegrpc.WithEndpoint(config.CollectorURL),
	)
	if err != nil {
		logrus.Fatalf("Failed to create trace exporter: %v", err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(otelResource()),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(provider)
	return exporter.Shutdown
}

func InitLogger() func(context.Context) error {
	var opt otlploggrpc.Option
	if config.Insecure == "false" {
		opt = otlploggrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	} else {
		opt = otlploggrpc.WithInsecure()
	}

	exporter, err := otlploggrpc.New(context.Background(),
		opt,
		otlploggrpc.WithEndpoint(config.CollectorURL),
	)
	if err != nil {
		logrus.Fatalf("Failed to create log exporter: %v", err)
	}

	provider := otel_log.NewLoggerProvider(
		otel_log.WithResource(otelResource()),
		otel_log.WithProcessor(otel_log.NewBatchProcessor(exporter)),
	)

	logrus.AddHook(otellogrus.NewHook(config.ServiceName, otellogrus.WithLoggerProvider(provider)))
	return provider.Shutdown
}

func InitMeter() func(context.Context) error {
	ctx := context.Background()
	var opt otlpmetricgrpc.Option
	if config.Insecure == "false" {
		opt = otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	} else {
		opt = otlpmetricgrpc.WithInsecure()
	}

	exporter, err := otlpmetricgrpc.New(ctx,
		opt,
		otlpmetricgrpc.WithEndpoint(config.CollectorURL),
	)
	if err != nil {
		logrus.Fatalf("Failed to create metric exporter: %v", err)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(otelResource()),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(10*time.Second))),
	)
	otel.SetMeterProvider(meterProvider)

	meter := meterProvider.Meter(config.ServiceName)
	OrdersTotal, _ = meter.Int64Counter("orders_total")
	ProductOrderCounter, _ = meter.Int64Counter("product_order_total")
	HttpRequestCount, _ = meter.Int64Counter("http_request_count")
	HttpDurationBucket, _ = meter.Float64Histogram("http_request_duration")
	GoroutinesGauge, _ = meter.Int64ObservableGauge("go_goroutines")
	MemoryGauge, _ = meter.Int64ObservableGauge("go_memory_bytes")

	meter.RegisterCallback(func(ctx context.Context, o metric.Observer) error {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		o.ObserveInt64(MemoryGauge, int64(memStats.Alloc))
		o.ObserveInt64(GoroutinesGauge, int64(runtime.NumGoroutine()))
		return nil
	}, MemoryGauge, GoroutinesGauge)

	otelruntime.Start()
	return meterProvider.Shutdown
}
