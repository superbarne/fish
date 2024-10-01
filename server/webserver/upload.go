package webserver

import (
	"errors"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/superbarne/fish/aquarium"
	"github.com/superbarne/fish/imageprocess"
	"github.com/superbarne/fish/models"
)

func (ws *WebServer) UploadFish(c *fiber.Ctx) error {
	// validate id
	aquariumID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Redirect("/aquarium")
	}

	var aquarium *aquarium.Aquarium
	var ok bool
	if err := func() error {
		ws.aquariumsLock.RLock()
		defer ws.aquariumsLock.RUnlock()
		aquarium, ok = ws.aquariums[aquariumID]
		if !ok {
			return errors.New("aquarium not found")
		}
		return nil
	}(); err != nil {
		return c.Redirect("/aquarium", fiber.StatusSeeOther)
	}

	if c.Method() == fiber.MethodPost {
		// Get the file from the request
		file, err := c.FormFile("image")
		if err != nil {
			ws.log.Error("Failed to get image from request", slog.String("error", err.Error()))
			return c.Redirect("/aquarium/"+aquariumID.String(), fiber.StatusSeeOther)
		}

		// is file a image
		if file.Header.Get("Content-Type") != "image/png" && file.Header.Get("Content-Type") != "image/jpeg" && file.Header.Get("Content-Type") != "image/jpg" {
			ws.log.Error("File is not a image", slog.String("content-type", file.Header.Get("Content-Type")))
			return c.Redirect("/aquarium/"+aquariumID.String(), fiber.StatusSeeOther)
		}

		name := c.FormValue("name", "Boid")
		fishID := uuid.New()

		tmpFilePath, err := ws.storage.SaveTmpFishImageFromRequest(c, aquarium.ID, fishID, file)
		if err != nil {
			ws.log.Error("Failed to save image", slog.String("error", err.Error()))
			return c.Redirect("/aquarium/"+aquariumID.String(), fiber.StatusSeeOther)
		}

		// Process Image
		targetPath := ws.storage.FishImagePath(aquarium.ID, fishID)
		if err := imageprocess.ProcessImage(tmpFilePath, targetPath, ws.log); err != nil {
			ws.log.Error("Failed to process image", slog.String("error", err.Error()))
			return c.Redirect("/aquarium/"+aquariumID.String(), fiber.StatusSeeOther)
		}

		// Remove tmp file
		if err := os.Remove(tmpFilePath); err != nil {
			ws.log.Error("Failed to remove temp file", slog.String("error", err.Error()))
		}

		// Write Json with metadata about the uploaded file
		fish := &models.Fish{
			ID:         fishID,
			AquariumID: aquarium.ID,
			Name:       name,
			Filename:   fishID.String() + ".png",
			Approved:   false,
		}

		if err := ws.storage.SaveFishMetadata(aquariumID, fish); err != nil {
			ws.log.Error("Failed to save fish metadata", slog.String("error", err.Error()))
			return c.Redirect("/aquarium/"+aquariumID.String(), fiber.StatusSeeOther)
		}

		aquarium.AddFish(fish)

		return c.Redirect("/aquarium/"+aquariumID.String(), fiber.StatusSeeOther)
	}

	return c.Render("upload", fiber.Map{
		"ID":    aquariumID.String(),
		"Title": "Go Fiber Template Example",
	})
}
