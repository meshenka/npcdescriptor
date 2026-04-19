package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/meshenka/npcgenerator/cmd"
	"github.com/meshenka/npcgenerator/internal/api"
	"github.com/meshenka/npcgenerator/internal/app"
	"github.com/meshenka/npcgenerator/internal/logging"
)

// @title NPC Descriptor API
// @version 1.0
// @description API for generating random NPC descriptors.
// @host localhost:8080
// @BasePath /
func main() {
	// Initialize structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := run(); err != nil {
		slog.Error("Fatal error", logging.Err(err))
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	mux := http.NewServeMux()
	application := app.NewApplication()
	api.RegisterRoutes(mux, application)

	// Serve static files from /public
	fs := http.FileServer(http.Dir("public"))
	mux.Handle("/", fs)

	port := cmd.GetEnv("PORT", "8080")
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           api.MuxMiddleware(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("Starting server", slog.String("url", "http://localhost:"+port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		slog.Info("Shutting down server...")
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	return srv.Shutdown(shutdownCtx)
}
