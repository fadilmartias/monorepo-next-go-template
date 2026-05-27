package middleware

import (
	"context"

	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func TraceIDMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {

		// Cek apakah client sudah kirim trace id
		traceID := c.Get("X-Request-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// Inject ke context
		ctx := c.Context()
		ctx = context.WithValue(ctx, utils.CtxTraceID, traceID)
		c.SetContext(ctx)

		// Set ke header response
		c.Set("X-Trace-ID", traceID)

		return c.Next()
	}
}

// fiber:context-methods migrated
