package client

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type FiberRedisStorage struct {
	client *redis.Client
	prefix string
}

func NewFiberRedisStorage(client *redis.Client, namespace string) *FiberRedisStorage {
	if namespace == "" {
		namespace = "fiber-shared-state"
	}

	return &FiberRedisStorage{
		client: client,
		prefix: namespace + ":",
	}
}

func (s *FiberRedisStorage) SetClient(client *redis.Client) {
	if s != nil {
		s.client = client
	}
}

func (s *FiberRedisStorage) GetWithContext(ctx context.Context, key string) ([]byte, error) {
	if s == nil || s.client == nil || key == "" {
		return nil, nil
	}

	val, err := s.client.Get(ctx, s.key(key)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (s *FiberRedisStorage) Get(key string) ([]byte, error) {
	return s.GetWithContext(context.Background(), key)
}

func (s *FiberRedisStorage) SetWithContext(ctx context.Context, key string, val []byte, exp time.Duration) error {
	if s == nil || s.client == nil || key == "" || len(val) == 0 {
		return nil
	}

	return s.client.Set(ctx, s.key(key), val, exp).Err()
}

func (s *FiberRedisStorage) Set(key string, val []byte, exp time.Duration) error {
	return s.SetWithContext(context.Background(), key, val, exp)
}

func (s *FiberRedisStorage) DeleteWithContext(ctx context.Context, key string) error {
	if s == nil || s.client == nil || key == "" {
		return nil
	}

	return s.client.Del(ctx, s.key(key)).Err()
}

func (s *FiberRedisStorage) Delete(key string) error {
	return s.DeleteWithContext(context.Background(), key)
}

func (s *FiberRedisStorage) ResetWithContext(ctx context.Context) error {
	if s == nil || s.client == nil {
		return nil
	}

	pattern := s.prefix + "*"
	iter := s.client.Scan(ctx, 0, pattern, 500).Iterator()
	keys := make([]string, 0, 500)

	flush := func() error {
		if len(keys) == 0 {
			return nil
		}
		if err := s.client.Unlink(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("fiber shared state reset failed: %w", err)
		}
		keys = keys[:0]
		return nil
	}

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
		if len(keys) >= cap(keys) {
			if err := flush(); err != nil {
				return err
			}
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}

	return flush()
}

func (s *FiberRedisStorage) Reset() error {
	return s.ResetWithContext(context.Background())
}

func (s *FiberRedisStorage) Close() error {
	return nil
}

func (s *FiberRedisStorage) key(key string) string {
	return s.prefix + strings.TrimPrefix(key, ":")
}
