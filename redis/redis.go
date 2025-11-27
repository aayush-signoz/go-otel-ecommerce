package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func SetupRedis(addr string) {
	RDB = redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

// Helper to set last order
func SetLastOrder(ctx context.Context, productID int, quantity int) {
	RDB.Set(ctx,
		fmt.Sprintf("last_order:%d", productID),
		quantity, 0,
	)
}
