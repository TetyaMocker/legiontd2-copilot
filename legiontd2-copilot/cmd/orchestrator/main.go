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
	"github.com/yourname/legiontd2-copilot/internal/perceptionclient"
	"github.com/yourname/legiontd2-copilot/internal/storage"
	"github.com/yourname/legiontd2-copilot/internal/webserver"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	slog.Info("lt2-copilot starting")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	dbPath := env("LT2_DB_PATH", "lt2_copilot.db")
	perceptionAddr := env("LT2_PERCEPTION_ADDR", "localhost:50051")
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

	percClient, err := perceptionclient.New(perceptionAddr)
	if err != nil {
		slog.Warn("perception service not available, running in offline mode", "error", err)
	}
	if percClient != nil {
		defer percClient.Close()
	}

	adv := advisor.NewHeuristicAdvisor()
	state := &webserver.AppState{}
	webserver.Start(state, webAddr)

	slog.Info("web UI available", "url", "http://localhost"+webAddr)

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	slog.Info("orchestrator running")
	for {
		select {
		case <-ctx.Done():
			slog.Info("shutting down")
			return
		case <-ticker.C:
			if percClient == nil {
				continue
			}

			readCtx, readCancel := context.WithTimeout(context.Background(), 16*time.Second)
			eco, err := percClient.ReadEconomy(readCtx)
			readCancel()
			if err != nil {
				slog.Warn("read economy failed", "error", err)
				continue
			}

			snap := advisor.EconomySnapshot{
				Mythium:           eco.Mythium,
				Income:            eco.Income,
				WaveNumber:        eco.WaveNumber,
				WaveTimerSeconds:  eco.WaveTimerSec,
				KingHPPercent:     eco.KingHPPercent,
				AllyKingHPPercent: eco.AllyKingHP,
				Confidence:        eco.Confidence,
			}

			recs := adv.Recommend(snap)
			state.Update(snap, recs)
		}
	}
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
