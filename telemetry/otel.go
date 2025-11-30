package telemetry

import (
	"context"
	"os"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/resource"

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
