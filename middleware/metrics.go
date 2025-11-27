package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var (
	HttpRequestCount   metric.Int64Counter
	HttpDurationBucket metric.Float64Histogram
)

// InitMetrics initializes counters using the global MeterProvider
func InitMetrics() error {
	meter := otel.Meter("go-otel-ecommerce")

	var err error
	HttpRequestCount, err = meter.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total number of HTTP requests"),
	)
	if err != nil {
		return err
	}

	HttpDurationBucket, err = meter.Float64Histogram(
		"http_request_duration_seconds",
		metric.WithDescription("Duration of HTTP requests in seconds"),
	)
	if err != nil {
		return err
	}

	return nil
}

// MetricsMiddleware records request count and duration
func MetricsMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next()
	duration := time.Since(start).Seconds()

	if HttpDurationBucket != nil {
		HttpDurationBucket.Record(c.Request.Context(), duration)
	}
	if HttpRequestCount != nil {
		HttpRequestCount.Add(c.Request.Context(), 1)
	}
}
