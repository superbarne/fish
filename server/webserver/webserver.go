package webserver

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/superbarne/fish/pubsub"
	"github.com/superbarne/fish/storage"
)

type WebServer struct {
	server *http.Server
	router chi.Router

	tmpl *template.Template
	log  *slog.Logger

	pubsub  *pubsub.PubSub
	storage *storage.Storage
}

func NewWebServer(log *slog.Logger, pubsub *pubsub.PubSub, store *storage.Storage) *WebServer {
	tmpl, err := template.ParseGlob("./views/*.html")
	if err != nil {
		log.Error("Failed to parse templates", slog.String("error", err.Error()))
		os.Exit(1)
	}

	ws := &WebServer{
		router:  chi.NewRouter(),
		tmpl:    tmpl,
		log:     log,
		pubsub:  pubsub,
		storage: store,
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
			r.Get("/fishes/{fishID}.png", ws.getFishImage)

			r.Group(func(r chi.Router) {
				r.Use(middleware.NoCache)

				r.Get("/", ws.uploadAquariumFish)
				r.Post("/", ws.uploadAquariumFish)
				r.Get("/sse", ws.sseAquarium)
			})
		})
		r.Handle("/*", http.StripPrefix("/aquarium", http.FileServer(http.Dir("./assets/aquarium"))))
	})

	ws.router.Route("/admin", func(r chi.Router) {
		r.Get("/", ws.listAdminAquariums)
		r.Route("/aquarium/{aquariumID:[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89aAbB][a-f0-9]{3}-[a-f0-9]{12}}", func(r chi.Router) {
			r.Get("/", ws.showAdminAquarium)
			r.Post("/approval", ws.toggleAdminNeedApproval)
			r.Route("/fishes/{fishID:[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89aAbB][a-f0-9]{3}-[a-f0-9]{12}}", func(r chi.Router) {
				r.Post("/delete", ws.deleteAdminFish)
				r.Post("/approve", ws.approveAdminFish)
			})
		})
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
