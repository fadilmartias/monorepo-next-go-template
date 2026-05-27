package client

import (
	"context"

	"github.com/fadilmartias/dilz_code/apps/backend/config"
	"github.com/go-redis/redis/v8"
)

func ConnectRedis() (*redis.Client, error) {
	redisConfig := config.LoadRedisConfig()
	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Host + ":" + redisConfig.Port,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

func Key(k string) string {
	redisConfig := config.LoadRedisConfig()
	if redisConfig == nil || redisConfig.Prefix == "" {
		return k
	}
	return redisConfig.Prefix + ":" + k
}

func DeleteKeysByPrefix(ctx context.Context, client *redis.Client, prefix string) error {
	batchSize := 500
	keysToDelete := make([]string, 0, batchSize)

	iter := client.Scan(ctx, 0, prefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		keysToDelete = append(keysToDelete, iter.Val())
		if len(keysToDelete) >= batchSize {
			if err := client.Unlink(ctx, keysToDelete...).Err(); err != nil {
				return err
			}
			keysToDelete = keysToDelete[:0]
		}
	}

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
