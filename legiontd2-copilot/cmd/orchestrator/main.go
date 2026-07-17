package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourname/legiontd2-copilot/internal/advisor"
	"github.com/yourname/legiontd2-copilot/internal/api"
	"github.com/yourname/legiontd2-copilot/internal/http"
	"github.com/yourname/legiontd2-copilot/internal/storage"
	"github.com/yourname/legiontd2-copilot/internal/ws"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	slog.Info("lt2-copilot starting")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	dbPath := env("LT2_DB_PATH", "lt2_copilot.db")
	webAddr := env("LT2_WEB_ADDR", ":8080")
	apiKey := os.Getenv("LT2_API_KEY")

	store, err := storage.New(dbPath)
	if err != nil {
		slog.Error("storage init failed", "error", err)
		os.Exit(1)
	}
	defer store.Close()

	if apiKey != "" {
		apiClient := api.NewClient(apiKey)
		_ = apiClient
		slog.Info("api client configured")
	}

	hub := ws.NewHub()
	httpserver.New(webAddr, hub)

	slog.Info("server ready", "url", "http://localhost"+webAddr)

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("shutting down")
			return
		case <-ticker.C:
			state := hub.GetState()
			recs := advisor.Recommend(state)
			hub.SetRecs(recs)
		}
	}
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
