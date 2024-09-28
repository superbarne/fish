package webserver

import (
	"context"
	"log/slog"
	"os"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/google/uuid"
	"github.com/superbarne/fish/aquarium"
)

type WebServer struct {
	app *fiber.App
	log *slog.Logger

	aquariums     map[uuid.UUID]*aquarium.Aquarium
	aquariumsLock sync.RWMutex
}

func NewWebServer(log *slog.Logger) *WebServer {
	// Initialize standard Go html template engine
	engine := html.New("./views", ".html")
	engine.Reload(true)

	// create folders
	os.MkdirAll("./uploads", os.ModePerm)
	os.MkdirAll("./data", os.ModePerm)

	ws := &WebServer{
		app: fiber.New(fiber.Config{
			Views: engine,
		}),
		log: log,
		aquariums: map[uuid.UUID]*aquarium.Aquarium{
			uuid.MustParse("38d7976d-3c27-4e74-8bfe-a9ec44318d3f"): aquarium.NewAquarium(uuid.MustParse("38d7976d-3c27-4e74-8bfe-a9ec44318d3f")),
		},
	}

	ws.app.Use(recover.New())
	ws.app.Use(compress.New())
	ws.app.Use(cors.New(
		cors.Config{
			AllowOrigins: "*",
		},
	))

	ws.app.Get("/", ws.ServeLandingPage)
	ws.app.Static("/assets", "./assets/app")
	ws.app.Static("/aquarium", "./assets/aquarium")
	ws.app.Static("/fishs", "./uploads")
	ws.app.Get("/aquarium/:id/sse", ws.ServeSSE)
	ws.app.Get("/aquarium/:id", ws.UploadFish)
	ws.app.Post("/aquarium/:id", ws.UploadFish)

	return ws
}

func (ws *WebServer) Listen() error {
	return ws.app.Listen(":8080")
}

func (ws *WebServer) Shutdown(ctx context.Context) error {
	return ws.app.ShutdownWithContext(ctx)
}
