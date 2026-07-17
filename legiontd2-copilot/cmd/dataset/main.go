package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/yourname/legiontd2-copilot/internal/api"
	"github.com/yourname/legiontd2-copilot/internal/dataset"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	apiKey := os.Getenv("LT2_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "LT2_API_KEY is required")
		os.Exit(1)
	}

	outDir := "dataset"
	if len(os.Args) > 1 {
		outDir = os.Args[1]
	}

	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir %s: %v\n", outDir, err)
		os.Exit(1)
	}

	client := api.NewClient(apiKey)
	collector := dataset.NewCollector(client, outDir)
	ctx := context.Background()

	// 1. Collect top 200 players by Elo
	slog.Info("step 1: collecting top players")
	topPlayers, err := collector.CollectTopPlayers(ctx, "overallElo", 200)
	if err != nil {
		slog.Error("collect top players", "error", err)
	}
	collector.SaveJSON("top_players", topPlayers)

	// 2. For each top player, collect match history
	type PlayerMatches struct {
		Player  api.Stats    `json:"player"`
		Matches []api.Match  `json:"matches"`
	}

	var allData []PlayerMatches
	playerLimit := 20
	if len(topPlayers) < playerLimit {
		playerLimit = len(topPlayers)
	}

	slog.Info("step 2: collecting match histories", "players", playerLimit)
	for i, p := range topPlayers[:playerLimit] {
		slog.Info("fetching matches", "player", p.PlayerName, "i", i+1, "total", playerLimit)

		matches, err := collector.CollectPlayerMatches(ctx, p.PlayerID, 20)
		if err != nil {
			slog.Warn("skip player", "player", p.PlayerName, "error", err)
			continue
		}

		allData = append(allData, PlayerMatches{
			Player:  p,
			Matches: matches,
		})

		time.Sleep(200 * time.Millisecond)
	}

	// 3. Save complete dataset
	collector.SaveJSON("dataset", allData)

	// 4. Save a compact version with just match details
	type CompactMatch struct {
		PlayerName string              `json:"playerName"`
		PlayerElo  int                 `json:"playerElo"`
		Version    string              `json:"version"`
		QueueType  string              `json:"queueType"`
		EndingWave int                 `json:"endingWave"`
		GameLength int                 `json:"gameLength"`
		GameResult string              `json:"gameResult"`
		Legion     string              `json:"legion"`
		Fighters   []string            `json:"fighters"`
		Mercs      []string            `json:"mercs"`
		WorkersPerWave  []int          `json:"workersPerWave"`
		IncomePerWave   []int          `json:"incomePerWave"`
		NetWorthPerWave []int          `json:"netWorthPerWave"`
		BuildPerWave    []string       `json:"buildPerWave"`
		LeaksPerWave    []string       `json:"leaksPerWave"`
		MercsSentPerWave []any         `json:"mercsSentPerWave"`
		KingUpgrades    []string       `json:"kingUpgrades"`
	}

	var compact []CompactMatch
	for _, pd := range allData {
		for _, m := range pd.Matches {
			for _, pl := range m.PlayersData {
				cm := CompactMatch{
					PlayerName: pl.PlayerName,
					PlayerElo:  pl.OverallElo,
					Version:    m.Version,
					QueueType:  m.QueueType,
					EndingWave: m.EndingWave,
					GameLength: m.GameLength,
					GameResult: pl.GameResult,
					Legion:     pl.Legion,
					WorkersPerWave:  pl.WorkersPerWave,
					IncomePerWave:   pl.IncomePerWave,
					NetWorthPerWave: pl.NetWorthPerWave,
					BuildPerWave:    pl.BuildPerWave,
					LeaksPerWave:    pl.LeaksPerWave,
					MercsSentPerWave: pl.MercenariesSentPerWave,
					KingUpgrades:    pl.KingUpgradesPerWave,
				}
				// Parse fighter/merc CSV
				if pl.Fighters != "" {
					cm.Fighters = splitCSV(pl.Fighters)
				}
				if pl.Mercenaries != "" {
					cm.Mercs = splitCSV(pl.Mercenaries)
				}
				compact = append(compact, cm)
			}
		}
	}

	collector.SaveJSON("dataset_compact", compact)

	// 5. Summary
	f, _ := os.Create(outDir + "/summary.txt")
	defer f.Close()
	fmt.Fprintf(f, "Dataset Summary\n")
	fmt.Fprintf(f, "==============\n\n")
	fmt.Fprintf(f, "Top players collected: %d\n", len(topPlayers))
	fmt.Fprintf(f, "Players with matches: %d\n", len(allData))
	fmt.Fprintf(f, "Total matches: %d\n", len(compact))
	fmt.Fprintf(f, "Total player match entries: %d\n", len(compact))

	var wins, losses int
	for _, m := range compact {
		if m.GameResult == "won" {
			wins++
		} else if m.GameResult == "lost" {
			losses++
		}
	}
	fmt.Fprintf(f, "Wins: %d, Losses: %d\n", wins, losses)

	slog.Info("done", "output", outDir)
	fmt.Printf("Dataset saved to %s/\n", outDir)
	fmt.Printf("  top_players.json — %d players\n", len(topPlayers))
	fmt.Printf("  dataset.json — full data\n")
	fmt.Printf("  dataset_compact.json — %d match entries\n", len(compact))
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	// The API uses comma-separated values, may have brackets for per-wave data
	var result []string
	current := ""
	for _, c := range s {
		if c == ',' {
			if current != "" {
				result = append(result, current)
			}
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}


