package model

import (
	"context"
	"errors"
	"ismismcube-backend/internal/config"
	"time"
)

func CheckAndSetRateLimit(clientIP string) error {
	ctx := context.Background()
	key := "danmu:rate_limit:" + clientIP
	result, err := config.RedisClient.SetNX(ctx, key, "1", time.Second).Result()
	if err != nil {
		return err
	}
	if !result {
		return errors.New("rate limit exceeded: only 1 danmu per second allowed")
	}
	return nil
}
