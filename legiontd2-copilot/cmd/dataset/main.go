package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/yourname/legiontd2-copilot/internal/api"
	"github.com/yourname/legiontd2-copilot/internal/dataset"
)

type CompactMatch struct {
	PlayerName string   `json:"playerName"`
	PlayerElo  int      `json:"playerElo"`
	Version    string   `json:"version"`
	QueueType  string   `json:"queueType"`
	EndingWave int      `json:"endingWave"`
	GameLength int      `json:"gameLength"`
	GameResult string   `json:"gameResult"`
	Legion     string   `json:"legion"`
	EloChange  int      `json:"eloChange"`
	Fighters   []string `json:"fighters"`
	Mercs      []string `json:"mercs"`
	Workers    []float64 `json:"workers"`
	Income     []float64 `json:"income"`
	NetWorth   []float64 `json:"netWorth"`
	Value      []float64 `json:"value"`
	ChosenSpell string  `json:"chosenSpell"`
	MvpScore   int      `json:"mvpScore"`
	LeakValue  int      `json:"leakValue"`
	Doubledown bool     `json:"doubledown"`
}

func parseFloatsField(raw json.RawMessage) []float64 {
	if len(raw) == 0 {
		return nil
	}
	// try array of numbers first
	var arr []float64
	if err := json.Unmarshal(raw, &arr); err == nil {
		return arr
	}
	// try space-separated string
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		if s == "" || s == "   " {
			return nil
		}
		parts := strings.Fields(s)
		out := make([]float64, 0, len(parts))
		for _, p := range parts {
			var v float64
			if _, err := fmt.Sscanf(p, "%f", &v); err == nil {
				out = append(out, v)
			}
		}
		return out
	}
	return nil
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
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

	gamesPath := filepath.Join(outDir, "games.json")
	var games []api.Match

	if _, err := os.Stat(gamesPath); err == nil {
		slog.Info("step 1: loading cached games")
		data, _ := os.ReadFile(gamesPath)
		json.Unmarshal(data, &games)
	}

	if len(games) == 0 {
		slog.Info("step 1: collecting games with details")
		var err error
		games, err = collector.CollectGames(ctx, "v26.6", 3000)
		if err != nil {
			slog.Error("collect games", "error", err)
		}
		collector.SaveJSON("games", games)
	} else {
		slog.Info("using cached games", "count", len(games))
	}

	slog.Info("step 2: flattening player data")
	var compact []CompactMatch
	for _, g := range games {
		var players []api.PlayerMatchDetails
		if len(g.PlayersData) > 0 {
			if err := json.Unmarshal(g.PlayersData, &players); err != nil {
				continue
			}
		}
		for _, pl := range players {
			cm := CompactMatch{
				PlayerName:  pl.PlayerName,
				PlayerElo:   pl.OverallElo,
				Version:     g.Version,
				QueueType:   g.QueueType,
				EndingWave:  g.EndingWave,
				GameLength:  g.GameLength,
				GameResult:  pl.GameResult,
				Legion:      pl.Legion,
				EloChange:   pl.EloChange,
				ChosenSpell: pl.ChosenSpell,
				MvpScore:    pl.MvpScore,
				LeakValue:   pl.LeakValue,
				Doubledown:  pl.Doubledown,
				Fighters:    splitCSV(pl.Fighters),
				Mercs:       splitCSV(pl.Mercenaries),
				Workers:     parseFloatsField(pl.WorkersPerWave),
				Income:      parseFloatsField(pl.IncomePerWave),
				NetWorth:    parseFloatsField(pl.NetWorthPerWave),
				Value:       parseFloatsField(pl.ValuePerWave),
			}
			compact = append(compact, cm)
		}
	}

	collector.SaveJSON("dataset_compact", compact)

	var wins, losses int
	for _, m := range compact {
		if m.GameResult == "won" {
			wins++
		} else if m.GameResult == "lost" {
			losses++
		}
	}

	f, _ := os.Create(outDir + "/summary.txt")
	defer f.Close()
	fmt.Fprintf(f, "Dataset Summary\n")
	fmt.Fprintf(f, "==============\n\n")
	fmt.Fprintf(f, "Games collected: %d\n", len(games))
	fmt.Fprintf(f, "Player entries: %d\n", len(compact))
	fmt.Fprintf(f, "Wins: %d, Losses: %d\n", wins, losses)
	fmt.Fprintf(f, "Win rate: %.1f%%\n", float64(wins)/float64(wins+losses)*100)

	slog.Info("done", "output", outDir, "games", len(games), "entries", len(compact))
	fmt.Printf("Dataset saved to %s/\n", outDir)
	fmt.Printf("  games.json — %d games\n", len(games))
	fmt.Printf("  dataset_compact.json — %d player entries\n", len(compact))
}
