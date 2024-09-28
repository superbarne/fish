package webserver

import (
	"bufio"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (ws *WebServer) ServeSSE(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		fmt.Fprintf(w, "event: ping\ndata: {}\n\n")
		w.Flush()

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fmt.Fprintf(w, "event: ping\ndata: {}\n\n")
				w.Flush()
			}
		}
	})

	return nil
}
