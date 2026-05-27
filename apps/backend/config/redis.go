package config

import (
	"os"
	"strconv"
	"sync"
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
