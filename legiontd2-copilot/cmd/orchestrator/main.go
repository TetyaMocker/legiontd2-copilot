package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourname/legiontd2-copilot/internal/advisor"
	"github.com/yourname/legiontd2-copilot/internal/deploy"
	"github.com/yourname/legiontd2-copilot/internal/http"
	"github.com/yourname/legiontd2-copilot/internal/storage"
	"github.com/yourname/legiontd2-copilot/internal/unitdata"
	"github.com/yourname/legiontd2-copilot/internal/ws"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	slog.Info("lt2-copilot starting")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	dbPath := env("LT2_DB_PATH", "lt2_copilot.db")
	webAddr := env("LT2_WEB_ADDR", ":8080")

	store, err := storage.New(dbPath)
	if err != nil {
		slog.Error("storage init failed", "error", err)
		os.Exit(1)
	}
	defer store.Close()

	if gameDir, err := deploy.FindGameDir(); err != nil {
		slog.Warn("game not found, skipping auto-deploy", "error", err)
	} else if err := deploy.DeployPatcher(gameDir); err != nil {
		slog.Error("auto-deploy failed", "error", err)
	} else {
		slog.Info("patcher ready", "gameDir", gameDir)
	}

	hub := ws.NewHub()

	iconsDir := ""
	if gameDir, err := deploy.FindGameDir(); err == nil {
		iconsDir = gameDir + `\Legion TD 2_Data\uiresources\AeonGT\hud\img\icons`
	}
	httpserver.New(webAddr, hub, iconsDir)

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
			for i := range state.Hand {
				supplementCost(&state.Hand[i])
			}
			for i := range state.Mercenaries {
				supplementCost(&state.Mercenaries[i])
			}
			for i := range state.TownActions {
				supplementCost(&state.TownActions[i])
			}
			recs := advisor.Recommend(state)
			hub.SetRecs(recs)
		}
	}
}

func supplementCost(u *ws.HandUnit) {
	if u.CostGold > 0 && u.CostSupply > 0 && u.CostMythium > 0 {
		return
	}
	if c, ok := unitdata.GetFighterCost(u.Name); ok {
		if u.CostGold == 0 {
			u.CostGold = c.Gold
		}
		if u.CostSupply == 0 {
			u.CostSupply = c.Supply
		}
		if u.CostMythium == 0 {
			u.CostMythium = c.Mythium
		}
	}
	if c, ok := unitdata.GetMercCost(u.Name); ok {
		if u.CostMythium == 0 {
			u.CostMythium = c.Mythium
		}
		if u.CostGold == 0 {
			u.CostGold = c.Gold
		}
		if u.CostSupply == 0 {
			u.CostSupply = c.Supply
		}
	}
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
