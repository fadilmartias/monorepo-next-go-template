package middleware

import (
	"strings"
	"time"

	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/golang-jwt/jwt"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/google/uuid"
)

func RateLimiter(max int, expiration time.Duration) fiber.Handler {
	if max == 0 {
		max = 100
	}
	if expiration == 0 {
		expiration = 1 * time.Minute
	}
	return limiter.New(limiter.Config{
		Next: func(c fiber.Ctx) bool {
			ip := c.Get("X-Forwarded-For")
			if ip == "" {
				ip = c.IP()
			}
			ip = strings.Split(ip, ",")[0]
			ip = strings.TrimSpace(ip)

			// Jangan limit kalau dari localhost
			return ip == "127.0.0.1" || ip == "::1"
		},
		Max:        max,
		Expiration: expiration,
		KeyGenerator: func(c fiber.Ctx) string {
			user := c.Locals("user")
			if user != nil {
				id := user.(jwt.MapClaims)["id"].(string)
				if id != "" {
					return id
				}
			}

			// cek visitor_id di cookie
			visitorID := c.Cookies("visitor_id")
			if visitorID == "" {
				visitorID = uuid.New().String()
				// set cookie biar persist
				c.Cookie(&fiber.Cookie{
					Name:     "visitor_id",
					Value:    visitorID,
					Expires:  time.Now().Add(365 * 24 * time.Hour), // 1 tahun
					HTTPOnly: true,
					SameSite: "Lax",
				})
			}
			return visitorID
		},
		LimitReached: func(c fiber.Ctx) error {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:    fiber.StatusTooManyRequests,
				Message: "Terlalu banyak permintaan",
			})
		},
		LimiterMiddleware: limiter.SlidingWindow{},
	})
}
