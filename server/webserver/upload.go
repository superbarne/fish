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
	"github.com/superbarne/fish/models"
)

func (ws *WebServer) UploadFish(c *fiber.Ctx) error {
	// validate id
	aquariumID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Redirect("/aquarium")
	}

	if c.Method() == fiber.MethodPost {
		// Get the file from the request
		file, err := c.FormFile("image")
		if err != nil {
			ws.log.Error("Failed to get image from request", slog.String("error", err.Error()))
			return c.Redirect("/aquarium/"+aquariumID.String(), fiber.StatusSeeOther)
		}

		name := c.FormValue("name", "Boid")

		// Generate a unique filename
		uniqueID := uuid.New()
		filename := fmt.Sprintf("%s%s", uniqueID.String(), filepath.Ext(file.Filename))

		// Save the file
		err = c.SaveFile(file, fmt.Sprintf("./uploads/%s", filename))
		if err != nil {
			ws.log.Error("Failed to save image", slog.String("error", err.Error()))
			return c.Redirect("/aquarium/"+aquariumID.String(), fiber.StatusSeeOther)
		}

		// Write Json with metadata about the uploaded file
		fish := &models.Fish{
			Name:       name,
			UploadTime: time.Now().String(),
			Approved:   false,
		}

		// Save Metadata to json file
		ws.saveToJSON(fish, uniqueID.String())

		return c.Redirect("/aquarium/"+aquariumID.String(), fiber.StatusSeeOther)
	}

	return c.Render("index", fiber.Map{
		"ID":    aquariumID.String(),
		"Title": "Go Fiber Template Example",
	})
}

func (ws *WebServer) saveToJSON(fish *models.Fish, uuid string) {
	// Convert the struct to JSON format
	jsonData, err := json.MarshalIndent(fish, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := os.MkdirAll("./data", os.ModePerm); err != nil {
		ws.log.Error("Failed to create data directory", slog.String("error", err.Error()))
		return
	}

	// Write the JSON data to a file
	file, err := os.Create(filepath.Join("./data", uuid+".json"))
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
