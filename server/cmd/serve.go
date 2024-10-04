package cmd

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/superbarne/fish/models"
	"github.com/superbarne/fish/pubsub"
	"github.com/superbarne/fish/storage"
	"github.com/superbarne/fish/webserver"
)

func NewServeCmd() *cobra.Command {
	serveCmd := &cobra.Command{
		Use: "serve",
		Run: func(cmd *cobra.Command, args []string) {
			serve()
		},
	}

	return serveCmd
}

func serve() {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	ps := pubsub.NewPubSub()
	store := storage.NewStorage("./data")

	// create default aquarium
	aquarium := &models.Aquarium{
		ID: uuid.MustParse("38d7976d-3c27-4e74-8bfe-a9ec44318d3f"),
	}
	if err := store.InsertAquarium(aquarium); err != nil {
		log.Error("Failed to insert default aquarium", slog.String("error", err.Error()))
		os.Exit(1)
		return
	}

	server := webserver.NewWebServer(log, ps, store)

	go func() {
		defer cancel()
		log.Info("Server start...")
		if err := server.Listen(); err != nil && err != http.ErrServerClosed {
			log.Error("Server error", slog.String("error", err.Error()))
		}
	}()

	<-ctx.Done()

	timeout, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	log.Info("Server shutdown...")

	server.Shutdown(timeout)
}
