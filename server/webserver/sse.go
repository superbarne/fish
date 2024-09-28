package webserver

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func (ws *WebServer) ServeSSE(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Status(fiber.StatusOK)

	c.Write([]byte("event: ping\ndata: ping\n\n"))

	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-c.Context().Done():
			return nil
		case <-ticker.C:
			// send ping
			c.Write([]byte("event: ping\ndata: ping\n\n"))
		}
	}
}
