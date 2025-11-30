package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"go-otel-ecommerce/config"
	"go-otel-ecommerce/db"
	"go-otel-ecommerce/handlers"
	"go-otel-ecommerce/middleware"
	"go-otel-ecommerce/redis"
	"go-otel-ecommerce/telemetry"
)

func main() {
	tracerShutdown := telemetry.InitTracer()
	defer tracerShutdown(context.Background())

	loggerShutdown := telemetry.InitLogger()
	defer loggerShutdown(context.Background())

	meterShutdown := telemetry.InitMeter()
	defer meterShutdown(context.Background())

	db.SetupDB()
	redis.SetupRedis(config.RedisAddr)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(otelgin.Middleware(config.ServiceName))
	r.Use(middleware.MetricsMiddleware)

	r.GET("/products", handlers.ListProducts)
	r.POST("/orders", handlers.CreateOrderHandler)
	r.GET("/checkInventory", handlers.CheckInventoryHandler)
	r.GET("/cpuTest", handlers.CpuLoadTest)
	r.GET("/concurrencyTest", handlers.ConcurrencyTest)

	logrus.WithField("service", config.ServiceName).Info(" Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
