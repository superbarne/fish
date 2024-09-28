package webserver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/superbarne/fish/aquarium"
)

func (ws *WebServer) ServeSSE(c *fiber.Ctx) error {
	// validate id
	aquariumID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Redirect("/aquarium")
	}

	var aquarium *aquarium.Aquarium
	func() {
		ws.aquariumsLock.RLock()
		defer ws.aquariumsLock.RUnlock()

		var ok bool
		if aquarium, ok = ws.aquariums[aquariumID]; !ok {
			c.Status(fiber.StatusNotFound)
			return
		}
	}()

	if aquarium == nil {
		c.Status(fiber.StatusNotFound)
		return nil
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	ctx := c.Context()
	fmt.Println(ctx)

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
			case fish := <-aquarium.Fishes(ctx):
				raw, _ := json.Marshal(fish)
				fmt.Fprintf(w, "event: fish\ndata: %s\n\n", raw)
			}
		}
	})

	return nil
}
