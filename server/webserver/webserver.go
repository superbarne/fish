package webserver

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type WebServer struct {
	app *fiber.App
	log *slog.Logger
}

func NewWebServer(log *slog.Logger) *WebServer {
	ws := &WebServer{
		app: fiber.New(),
		log: log,
	}

	ws.app.Use(recover.New())
	ws.app.Use(compress.New())

	ws.app.Static("/aquarium", "./assets/aquarium")
	ws.app.Get("/aquarium/:id/sse", ws.ServeSSE)

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
