package middleware

import (
	"context"
	"strings"

	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt"
)

func GetRealIP(c fiber.Ctx) string {
	// Prioritas 1: Cloudflare
	cfIP := c.Get("CF-Connecting-IP")
	if cfIP != "" {
		return cfIP
	}

	// Prioritas 2: X-Forwarded-For (ambil pertama)
	xff := c.Get("X-Forwarded-For")
	if xff != "" {
		// bisa berisi banyak IP, format: "ip1, ip2, ip3"
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}

	// Prioritas 3: X-Real-IP
	realIP := c.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// fallback: IP dari Fiber (biasanya Cloudflare IP)
	return c.IP()
}

// ActivityContextMiddleware injects causer & request metadata into the user context
func ActivityContextMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {

		// Ambil user ID dari c.Locals("user_id")
		// misal di auth middleware kamu sudah set ini
		user := c.Locals("user")
		var userID string

		if user != nil {
			userID = user.(jwt.MapClaims)["id"].(string)
		}

		ip := GetRealIP(c)

		// generate context baru
		ctx := c.Context()
		ctx = context.WithValue(ctx, utils.CtxCauserID, userID)
		ctx = context.WithValue(ctx, utils.CtxIP, ip)
		ctx = context.WithValue(ctx, utils.CtxUserAgent, string(c.Request().Header.UserAgent()))

		// assign context baru ke Fiber
		c.SetContext(ctx)

		return c.Next()
	}
}

// fiber:context-methods migrated
