package jobs

import (
	"context"
	"log"

	"github.com/fadilmartias/dilz_code/apps/backend/config"
	"github.com/hibiken/asynq"
)

var (
	AsynqClient *asynq.Client
	AsynqServer *asynq.Server
)

func InitQueue(redisClient *config.RedisClient) {
	// Pakai redis options yang sama
	AsynqClient = asynq.NewClient(asynq.RedisClientOpt{
		Addr:     redisClient.GetClient().Options().Addr,
		Password: redisClient.GetClient().Options().Password,
		DB:       redisClient.GetClient().Options().DB,
	})

	AsynqServer = asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     redisClient.GetClient().Options().Addr,
			Password: redisClient.GetClient().Options().Password,
			DB:       redisClient.GetClient().Options().DB,
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
