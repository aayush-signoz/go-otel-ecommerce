package handlers

import (
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"

	"go-otel-ecommerce/db"
	"go-otel-ecommerce/models"
	"go-otel-ecommerce/redis"
	"go-otel-ecommerce/telemetry"
)

// Handlers

func ListProducts(c *gin.Context) {
	rows, err := db.DB.QueryContext(c.Request.Context(), "SELECT id, name FROM products")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch products"})
		return
	}
	defer rows.Close()

	var products []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		products = append(products, map[string]interface{}{"id": id, "name": name})
	}
	c.JSON(http.StatusOK, gin.H{"products": products})
}

func CreateOrderHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var req models.CreateOrderRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tr := otel.Tracer("goApp")
	_, span := tr.Start(ctx, "create_order")
	defer span.End()

	var productID int
	err := db.DB.QueryRow("SELECT id FROM products WHERE name = ?", req.ProductName).Scan(&productID)
	if err != nil {
		res, err := db.DB.Exec("INSERT INTO products (name) VALUES (?)", req.ProductName)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create product"})
			return
		}
		lastID, _ := res.LastInsertId()
		productID = int(lastID)
	}

	_, err = db.DB.Exec("INSERT INTO orders (product_id, quantity, user_id) VALUES (?, ?, ?)",
		productID, req.Quantity, req.UserID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create order"})
		return
	}

	redis.SetLastOrder(ctx, productID, req.Quantity)

	telemetry.OrdersTotal.Add(ctx, int64(req.Quantity))
	telemetry.ProductOrderCounter.Add(ctx, int64(req.Quantity),
		metric.WithAttributes(attribute.String("product", req.ProductName)),
	)

	span.SetAttributes(
		attribute.String("product.name", req.ProductName),
		attribute.Int("quantity", req.Quantity),
	)

	logrus.WithContext(ctx).Infof("Order created | product=%s qty=%d", req.ProductName, req.Quantity)
	c.JSON(http.StatusOK, gin.H{"status": "order created"})
}

func CheckInventoryHandler(c *gin.Context) {
	ctx := c.Request.Context()
	delay := time.Duration(100+rand.Intn(300)) * time.Millisecond
	time.Sleep(delay)

	logrus.WithContext(ctx).Infof("Inventory checked | delay=%dms", delay.Milliseconds())
	c.JSON(http.StatusOK, gin.H{
		"inventory_status": "in stock",
		"check_time_ms":    delay.Milliseconds(),
	})
}

func CpuLoadTest(c *gin.Context) {
	start := time.Now()
	for i := 0; i < 9_000_000; i++ {
		_ = i * i
	}
	c.JSON(200, gin.H{"cpu_test_ms": time.Since(start).Milliseconds()})
}

func ConcurrencyTest(c *gin.Context) {
	var wg sync.WaitGroup
	for i := 0; i < 300; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(time.Millisecond * 200)
			wg.Done()
		}()
	}
	wg.Wait()
	c.JSON(200, gin.H{"goroutines": runtime.NumGoroutine()})
}
