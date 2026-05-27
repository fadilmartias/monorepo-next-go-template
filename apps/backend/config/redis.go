package config

import (
	"os"
	"strconv"
	"sync"
)

type RedisConfig struct {
	Host     string
	Port     string
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
			Host:     os.Getenv("REDIS_HOST"),
			Port:     os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       db,
			Prefix:   os.Getenv("REDIS_PREFIX"),
		}
	})
	return redisConfig
}
