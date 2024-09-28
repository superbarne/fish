package webserver

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

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
	func() {
		ws.aquariumsLock.RLock()
		defer ws.aquariumsLock.RUnlock()
		aquarium, ok = ws.aquariums[aquariumID]
		if !ok {
			c.Redirect("/aquarium")
			return
		}
	}()

	if c.Method() == fiber.MethodPost {
		// Get the file from the request
		file, err := c.FormFile("image")
		if err != nil {
			ws.log.Error("Failed to get image from request", slog.String("error", err.Error()))
			return c.Redirect("/aquarium/"+aquariumID.String(), fiber.StatusSeeOther)
		}

		name := c.FormValue("name", "Boid")

		// Generate a unique filename
		fishID := uuid.New()
		filename := fmt.Sprintf("%s%s", fishID.String(), filepath.Ext(file.Filename))

		tmpFile := filepath.Join(os.TempDir(), aquariumID.String(), filename)
		if err := os.MkdirAll(filepath.Join(os.TempDir(), aquariumID.String()), os.ModePerm); err != nil {
			ws.log.Error("Failed to create data directory", slog.String("error", err.Error()))
			return c.Redirect("/aquarium/"+aquariumID.String(), fiber.StatusSeeOther)
		}

		// Save the file
		err = c.SaveFile(file, tmpFile)
		if err != nil {
			ws.log.Error("Failed to save image", slog.String("error", err.Error()))
			return c.Redirect("/aquarium/"+aquariumID.String(), fiber.StatusSeeOther)
		}

		// Process Image
		targetFile := filepath.Join("./uploads/", aquariumID.String(), fishID.String()+".png")
		imageprocess.ProcessImage(tmpFile, targetFile)

		// Remove temp file
		if err := os.Remove(tmpFile); err != nil {
			ws.log.Error("Failed to remove temp file", slog.String("error", err.Error()))
		}

		// Write Json with metadata about the uploaded file
		fish := &models.Fish{
			ID:         fishID,
			Name:       name,
			Filename:   fishID.String() + ".png",
			UploadTime: time.Now().String(),
			Approved:   false,
		}

		// Save Metadata to json file
		ws.saveToJSON(aquarium.ID, fish, fishID.String())

		aquarium.AddFish(fish)

		return c.Redirect("/aquarium/"+aquariumID.String(), fiber.StatusSeeOther)
	}

	return c.Render("index", fiber.Map{
		"ID":    aquariumID.String(),
		"Title": "Go Fiber Template Example",
	})
}

func (ws *WebServer) saveToJSON(aquariumID uuid.UUID, fish *models.Fish, uuid string) {
	// Convert the struct to JSON format
	jsonData, err := json.MarshalIndent(fish, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := os.MkdirAll("./data/"+aquariumID.String(), os.ModePerm); err != nil {
		ws.log.Error("Failed to create data directory", slog.String("error", err.Error()))
		return
	}

	// Write the JSON data to a file
	file, err := os.Create(filepath.Join("./data", aquariumID.String(), uuid+".json"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println(err)
		return
	}
}
