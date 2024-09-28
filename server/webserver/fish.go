package webserver

import (
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (ws *WebServer) ServeFishImage(c *fiber.Ctx) error {
	// validate id
	aquariumID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Redirect("/aquarium")
	}

	fishID, err := uuid.Parse(c.Params("fishID"))
	if err != nil {
		return c.Redirect("/aquarium/" + aquariumID.String())
	}

	return c.SendFile(filepath.Join("./uploads", aquariumID.String(), fishID.String()))
}
