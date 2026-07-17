package dataset

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/yourname/legiontd2-copilot/internal/api"
)

type Collector struct {
	client *api.Client
	outDir string
	mu     sync.Mutex
}

func NewCollector(client *api.Client, outDir string) *Collector {
	return &Collector{client: client, outDir: outDir}
}

func (c *Collector) CollectTopPlayers(ctx context.Context, sortBy string, count int) ([]api.Stats, error) {
	slog.Info("collecting top players", "sortBy", sortBy, "count", count)
	var all []api.Stats
	offset := 0
	pageSize := 100
	for len(all) < count {
		n := pageSize
		if n > count-len(all) {
			n = count - len(all)
		}
		stats, err := c.client.GetTopPlayers(ctx, sortBy, n, offset)
		if err != nil {
			return all, fmt.Errorf("get top players at offset %d: %w", offset, err)
		}
		if len(stats) == 0 {
			break
		}
		all = append(all, stats...)
		offset += len(stats)
		time.Sleep(100 * time.Millisecond)
	}
	return all, nil
}

func (c *Collector) CollectPlayerMatches(ctx context.Context, playerID string, maxMatches int) ([]api.Match, error) {
	slog.Info("collecting matches", "player", playerID, "max", maxMatches)
	var all []api.Match
	offset := 0
	pageSize := 50
	for len(all) < maxMatches {
		n := pageSize
		if n > maxMatches-len(all) {
			n = maxMatches - len(all)
		}
		matches, err := c.client.GetMatchHistory(ctx, playerID, n, offset, true)
		if err != nil {
			return all, fmt.Errorf("get match history for %s at offset %d: %w", playerID, offset, err)
		}
		if len(matches) == 0 {
			break
		}
		all = append(all, matches...)
		offset += len(matches)
		time.Sleep(100 * time.Millisecond)
	}
	return all, nil
}

func (c *Collector) SaveJSON(name string, data any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := os.MkdirAll(c.outDir, 0755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	path := filepath.Join(c.outDir, name+".json")
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("encode %s: %w", path, err)
	}

	slog.Info("saved dataset", "path", path)
	return nil
}
