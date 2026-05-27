package utils

import (
	"github.com/fadilmartias/dilz_code/apps/backend/app/client"
	"github.com/gofiber/fiber/v3"
)

func Http() *client.HttpClient {
	return client.Http()
}

func HttpWithCtx(c fiber.Ctx) *client.HttpClient {
	return client.HttpWithCtx(c)
}
