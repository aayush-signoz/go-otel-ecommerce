package config

import "os"

var (
	ServiceName  = os.Getenv("SERVICE_NAME")
	CollectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	Insecure     = os.Getenv("INSECURE_MODE")
	RedisAddr    = os.Getenv("REDIS_ADDR")
)
