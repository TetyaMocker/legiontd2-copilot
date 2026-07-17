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
	"github.com/yourname/legiontd2-copilot/internal/storage"
	"github.com/yourname/legiontd2-copilot/internal/webserver"
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

	adv := advisor.NewHeuristicAdvisor()
	state := &webserver.AppState{}
	webserver.Start(state, webAddr)

	slog.Info("web UI + ingestion API", "url", "http://localhost"+webAddr)

	recTicker := time.NewTicker(2 * time.Second)
	defer recTicker.Stop()

	slog.Info("orchestrator running")
	for {
		select {
		case <-ctx.Done():
			slog.Info("shutting down")
			return
		case <-recTicker.C:
			snap := state.Snapshot()
			eco := advisor.EconomySnapshot{
				Mythium:          snap.Mythium,
				Income:           snap.Income,
				WaveNumber:       snap.Wave,
				WaveTimerSeconds: snap.WaveTimer,
				KingHPPercent:    snap.KingHP,
				Confidence:       float32(snap.Confidence),
			}
			recs := adv.Recommend(eco)
			state.SetRecs(recs)
		}
	}
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
