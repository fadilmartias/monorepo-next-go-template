package jobs

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
)

var (
	AsynqClient *asynq.Client
	AsynqServer *asynq.Server
)

func InitQueue(redisClient *redis.Client) {
	// Pakai redis options yang sama
	options := redisClient.Options()
	AsynqClient = asynq.NewClient(asynq.RedisClientOpt{
		Addr:     options.Addr,
		Password: options.Password,
		DB:       options.DB,
	})

	AsynqServer = asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     options.Addr,
			Password: options.Password,
			DB:       options.DB,
		},
		asynq.Config{
			Concurrency: 10,
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Printf("Task %s failed: %v", task.Type(), err)
			}),
			// Queues: map[string]int{
			// 	"critical": 10,
			// 	"high":     10,
			// 	"default":  10,
			// 	"low":      10,
			// },
		},
	)
}
