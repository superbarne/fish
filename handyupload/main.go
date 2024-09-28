package main

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/google/uuid"
)

type Fish struct {
	Name string `json:"name"`
	UploadTime string `json:"upload_time"`
	Approved bool `json:"approved"`
}

func main() {
    // Initialize standard Go html template engine
    engine := html.New("./views", ".html")

    app := fiber.New(fiber.Config{
        Views: engine,
    })

    app.Get("/", func(c *fiber.Ctx) error {
        // Render index template
        return c.Render("index", fiber.Map{
            "Title": "Go Fiber Template Example",
            "Description": "An example template",
            "Greeting": "Hello, world!",
        });
    })

	app.Post("/upload", handleImageUpload)
	app.Static("/images", "./uploads")

    log.Fatal(app.Listen(":3000"))
}

func handleImageUpload(c *fiber.Ctx) error {
    // Get the file from the request
    file, err := c.FormFile("image")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Failed to get image from request",
        })
    }

	name := c.FormValue("name", "Boid")

    // Generate a unique filename
    uniqueID := uuid.New()
    filename := fmt.Sprintf("%s%s", uniqueID.String(), filepath.Ext(file.Filename))

    // Save the file
    err = c.SaveFile(file, fmt.Sprintf("./uploads/%s", filename))
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to save image",
        })
    }

	// Write Json with metadata about the uploaded file
	fish := Fish{
		Name: name,
		UploadTime: time.Now().String(),
		Approved: false,
	}

	// Save Metadata to json file


    // Return success response
    return c.JSON(fiber.Map{
        "filename": filename,
        "size": file.Size,
    })
}
