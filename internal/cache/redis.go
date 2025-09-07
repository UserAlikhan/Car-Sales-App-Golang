package cache

import (
	"car_sales/internal/configs"
	"time"
)

func SetCache(key string, value string, ttl time.Duration) error {
	return configs.RedisClient.Set(configs.Ctx, key, value, ttl).Err()
}

func GetCache(key string) (string, error) {
	return configs.RedisClient.Get(configs.Ctx, key).Result()
}

func DeleteCache(key string) error {
	return configs.RedisClient.Del(configs.Ctx, key).Err()
}
