package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourname/legiontd2-copilot/internal/advisor"
	"github.com/yourname/legiontd2-copilot/internal/perceptionclient"
	"github.com/yourname/legiontd2-copilot/internal/storage"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	slog.Info("lt2-copilot starting")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	dbPath := os.Getenv("LT2_DB_PATH")
	if dbPath == "" {
		dbPath = "lt2_copilot.db"
	}

	perceptionAddr := os.Getenv("LT2_PERCEPTION_ADDR")
	if perceptionAddr == "" {
		perceptionAddr = "localhost:50051"
	}

	store, err := storage.New(dbPath)
	if err != nil {
		slog.Error("failed to init storage", "error", err)
		os.Exit(1)
	}
	defer store.Close()

	percClient, err := perceptionclient.New(perceptionAddr)
	if err != nil {
		slog.Warn("perception service not available", "error", err)
	}
	if percClient != nil {
		defer percClient.Close()
	}

	adv := advisor.NewHeuristicAdvisor()

	// основной цикл опроса Perception Service
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if percClient == nil {
					continue
				}

				eco, err := percClient.ReadEconomy(context.Background())
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

				slog.Info("advisory tick",
					"mythium", eco.Mythium,
					"income", eco.Income,
					"wave", eco.WaveNumber,
					"confidence", eco.Confidence,
					"recommendations", len(recs),
				)

				for _, r := range recs {
					slog.Info("recommendation", "kind", r.Kind, "explanation", r.Explanation)
				}
			}
		}
	}()

	slog.Info("orchestrator running")
	<-ctx.Done()
	slog.Info("shutting down")
}
