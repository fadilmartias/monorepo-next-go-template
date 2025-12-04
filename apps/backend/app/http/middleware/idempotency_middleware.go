package middleware

import (
	"bytes"
	"time"

	"github.com/bytedance/sonic"
	"github.com/fadilmartias/dilz_code/apps/backend/config"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type cachedResponse struct {
	Status int                    `json:"status"`
	Body   sonic.NoCopyRawMessage `json:"body"`
}

// Lama penyimpanan hasil (misal 5 menit)
const cacheTTL = 5 * time.Minute

func Idempotency(redis *config.RedisClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if fiber.IsMethodSafe(c.Method()) {
			return c.Next()
		}
		key := c.Get("X-Idempotency-Key")
		if key == "" {
			key = uuid.NewString()
			c.Request().Header.Add("X-Idempotency-Key", key)
		}

		ctx := c.UserContext()
		cacheKey := "idem:" + key

		// 🔹 Cek apakah sudah pernah disimpan
		val, err := redis.Get(ctx, cacheKey)
		if err == nil && val != "" {
			var cached cachedResponse
			if err := sonic.Unmarshal([]byte(val), &cached); err == nil {
				return c.Status(cached.Status).Send(cached.Body)
			}
		}

		// 🔹 Tangkap respons handler
		buf := new(bytes.Buffer)
		c.Response().SetBodyStream(buf, -1)

		if err := c.Next(); err != nil {
			return err
		}

		// 🔹 Ambil data hasil response
		status := c.Response().StatusCode()
		body := c.Response().Body()

		cached := cachedResponse{
			Status: status,
			Body:   sonic.NoCopyRawMessage(body),
		}

		jsonData, _ := sonic.Marshal(cached)
		_ = redis.Set(ctx, cacheKey, jsonData, cacheTTL)

		return nil
	}
}
