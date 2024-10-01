package cmd

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
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

	store := storage.NewStorage("./data")

	server := webserver.NewWebServer(log, store)

	go func() {
		defer cancel()
		log.Info("Server start...")
		if err := server.Listen(); err != nil {
			log.Error("Server error", slog.String("error", err.Error()))
		}
	}()

	<-ctx.Done()

	timeout, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	log.Info("Server shutdown...")

	server.Shutdown(timeout)
}
