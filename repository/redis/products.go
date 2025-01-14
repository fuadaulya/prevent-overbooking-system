package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// GetProductStockFromCache retrieves the stock of a product from Redis
func (r *Redis) GetProductStockFromCache(ctx context.Context, productID int) (int, error) {
	var stock int

	stockStr, err := r.client.Get(ctx, fmt.Sprintf("product_stock:%d", productID)).Result()
	if err != nil {
		if err == redis.Nil {
			// Key not found in Redis
			return stock, nil
		}
		return stock, fmt.Errorf("[redis] failed to get stock from Redis: %+v", err)
	}

	// Convert stock from string to int
	stock, err = strconv.Atoi(stockStr)
	if err != nil {
		return stock, fmt.Errorf("[redis] failed to parse stock: %+v", err)
	}

	return stock, nil
}

// SetProductStockInCache stores the stock of a product in Redis
func (r *Redis) SetProductStockInCache(ctx context.Context, productID, stock int) error {
	key := fmt.Sprintf("product_stock:%d", productID)

	_, err := r.client.Set(ctx, key, stock, 30*time.Second).Result()
	if err != nil {
		return fmt.Errorf("[redis] failed to set stock in Redis: %+v", err)
	}
	return nil
}
