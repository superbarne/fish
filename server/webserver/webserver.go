package webserver

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/superbarne/fish/aquarium"
	"github.com/superbarne/fish/storage"
)

type WebServer struct {
	server *http.Server
	router chi.Router

	tmpl    *template.Template
	log     *slog.Logger
	storage *storage.Storage

	aquariums     map[uuid.UUID]*aquarium.Aquarium
	aquariumsLock sync.RWMutex
}

func NewWebServer(log *slog.Logger, store *storage.Storage) *WebServer {
	tmpl, err := template.ParseGlob("./views/*.html")
	if err != nil {
		log.Error("Failed to parse templates", slog.String("error", err.Error()))
		os.Exit(1)
	}

	ws := &WebServer{
		router:  chi.NewRouter(),
		tmpl:    tmpl,
		log:     log,
		storage: store,
		aquariums: map[uuid.UUID]*aquarium.Aquarium{
			uuid.MustParse("38d7976d-3c27-4e74-8bfe-a9ec44318d3f"): aquarium.NewAquarium(uuid.MustParse("38d7976d-3c27-4e74-8bfe-a9ec44318d3f"), store),
		},
	}

	// add chi middlewares
	ws.router.Use(middleware.Recoverer)
	ws.router.Use(middleware.RequestID)
	ws.router.Use(middleware.Compress(5, "gzip"))
	ws.router.Use(middleware.StripSlashes)
	ws.router.Use(cors.AllowAll().Handler)

	ws.router.Get("/", ws.getLandingPage)
	ws.router.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/app"))))

	ws.router.Route("/aquarium", func(r chi.Router) {
		r.Route("/{aquariumID:[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89aAbB][a-f0-9]{3}-[a-f0-9]{12}}", func(r chi.Router) {
			r.Get("/fishs/{fishID}.png", ws.getFishImage)

			r.Group(func(r chi.Router) {
				r.Use(middleware.NoCache)
				r.Get("/", ws.uploadFish)
				r.Post("/", ws.uploadFish)
				r.Get("/sse", ws.sseFish)
			})
		})
		r.Handle("/*", http.StripPrefix("/aquarium", http.FileServer(http.Dir("./assets/aquarium"))))
	})

	return ws
}

func (ws *WebServer) Listen() error {
	// read env
	port := os.Getenv("AQUARIUM_PORT")
	if port == "" {
		port = "3000"
	}

	ws.server = &http.Server{
		Addr:    ":" + port,
		Handler: ws.router,
	}

	ws.log.Info("WebServer is running", slog.String("port", port))

	return ws.server.ListenAndServe()
}

func (ws *WebServer) Shutdown(ctx context.Context) error {
	return ws.server.Shutdown(ctx)
}
