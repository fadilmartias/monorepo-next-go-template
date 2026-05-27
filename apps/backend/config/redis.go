package config

import (
	"context"
	"os"
	"strconv"
	"sync"

	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	Prefix   string
}

var (
	redisConfig *RedisConfig
	redisOnce   sync.Once
)

func LoadRedisConfig() *RedisConfig {
	redisOnce.Do(func() {
		db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
		redisConfig = &RedisConfig{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       db,
			Prefix:   os.Getenv("REDIS_PREFIX"),
		}
	})
	return redisConfig
}

func Key(k string) string {
	redisConfig := LoadRedisConfig()
	if redisConfig == nil || redisConfig.Prefix == "" {
		return k
	}
	return redisConfig.Prefix + ":" + k
}

func DeleteKeysByPrefix(ctx context.Context, client *redis.Client, prefix string) error {
	batchSize := 500 // jumlah key dihapus per batch
	keysToDelete := make([]string, 0, batchSize)

	iter := client.Scan(ctx, 0, prefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		keysToDelete = append(keysToDelete, iter.Val())

		// Jika sudah mencapai batchSize, hapus dengan UNLINK biar non-blocking
		if len(keysToDelete) >= batchSize {
			if err := client.Unlink(ctx, keysToDelete...).Err(); err != nil {
				return err
			}
			keysToDelete = keysToDelete[:0]
		}
	}

	// Hapus sisa key kalau ada
	if len(keysToDelete) > 0 {
		if err := client.Unlink(ctx, keysToDelete...).Err(); err != nil {
			return err
		}
	}

	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}
