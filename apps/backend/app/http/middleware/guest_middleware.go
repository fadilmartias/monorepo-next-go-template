package middleware

import (
	"errors"
	"os"

	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/gofiber/fiber/v3"
)

func Guest() fiber.Handler {
	return func(c fiber.Ctx) error {
		if c.Value("user") != nil || c.Get("Authorization") != "" || c.Get("X-API-Key") != "" || c.Cookies("access_token_"+os.Getenv("APP_ENV")) != "" {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:    fiber.StatusForbidden,
				Message: "Forbidden: User already logged in",
			}, errors.New("user already logged in"))
		}
		return c.Next()
	}
}
