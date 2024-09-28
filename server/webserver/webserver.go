package webserver

import (
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
)

type WebServer struct {
	app *fiber.App
	log *slog.Logger
}

func NewWebServer(log *slog.Logger) *WebServer {
	// Initialize standard Go html template engine
	engine := html.New("./views", ".html")

	// create folders
	os.MkdirAll("./uploads", os.ModePerm)
	os.MkdirAll("./data", os.ModePerm)

	ws := &WebServer{
		app: fiber.New(fiber.Config{
			Views: engine,
		}),
		log: log,
	}

	ws.app.Use(recover.New())
	ws.app.Use(compress.New())

	ws.app.Static("/aquarium", "./assets/aquarium")
	ws.app.Get("/aquarium/:id/sse", ws.ServeSSE)
	ws.app.Get("/aquarium/:id", ws.UploadFish)
	ws.app.Post("/aquarium/:id", ws.UploadFish)

	// serve aquarium
	// serve upload
	// serve serve fishis

	return ws
}

func (ws *WebServer) Listen() error {
	return ws.app.Listen(":8080")
}

func (ws *WebServer) Shutdown() {
	ws.app.Shutdown()
}
