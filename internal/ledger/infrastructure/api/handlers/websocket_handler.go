package handlers

import (
	ws "lenavire/internal/ledger/infrastructure/websocket"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func HandleWebSocket(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func HandleLedgerActivity(hub *ws.LedgerActivityHub) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		client := &ws.Client{Conn: c}
		hub.Register <- client

		defer func() {
			hub.Unregister <- client
		}()

		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				break
			}
		}
	})
}
