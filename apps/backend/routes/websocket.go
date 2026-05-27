package routes

import (
	"fmt"

	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

func RegisterWebsocketRoutes(app *fiber.App) {

	app.Use("/ws", func(c fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/payment/:order_id", websocket.New(func(c *websocket.Conn) {
		fmt.Println("New websocket client connected: ", c.Params("order_id"))
		orderID := c.Params("order_id")
		utils.WebsocketAddClient("payment:"+orderID, c)
		defer utils.WebsocketRemoveClient("payment:"+orderID, c)

		for {
			// Kalau client kirim pesan, bisa dibaca di sini (opsional)
			if _, _, err := c.ReadMessage(); err != nil {
				break
			}
		}
	}))

	app.Get("/ws/chat/:room", websocket.New(func(c *websocket.Conn) {
		room := c.Params("room")
		utils.WebsocketAddClient("chat:"+room, c)
		defer utils.WebsocketRemoveClient("chat:"+room, c)
	}))

}
