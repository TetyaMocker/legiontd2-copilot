// Package api — клиент официального Legion TD 2 API v2 (https://apiv2.legiontd2.com).
//
// Используется ТОЛЬКО для офлайн-контура: справочники юнитов/волн/спеллов,
// история собственных матчей для обучения Advisor'а (см. ТЗ, раздел 8.4).
// Никогда не вызывается в горячем пути во время матча — не даёт live-данных.
//
// Ключ выпускается на developer.legiontd2.com, передаётся в заголовке x-api-key.
// Rate limit по состоянию на Phase 0: 5 req/s, 100 burst, 1000/день.
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const baseURL = "https://apiv2.legiontd2.com"

type Client struct {
	httpClient *http.Client
	apiKey     string
}

func NewClient(apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		apiKey:     apiKey,
	}
}

func (c *Client) doRequest(ctx context.Context, path string, dest any) error {
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+path, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("x-api-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("api returned status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

type Unit struct {
	ID          string `json:"_id"`
	Name        string `json:"name"`
	MythiumCost string `json:"mythiumCost"`
	GoldCost    string `json:"goldcost"`
	HP          string `json:"hp"`
	UnitClass   string `json:"unitClass"`
}

func (c *Client) GetUnitsByVersion(ctx context.Context, version string) ([]Unit, error) {
	var units []Unit
	if err := c.doRequest(ctx, "/units/byVersion/"+version, &units); err != nil {
		return nil, fmt.Errorf("get units: %w", err)
	}
	return units, nil
}

func (c *Client) GetMatchHistory(ctx context.Context, playerID string, limit, offset int) ([]Match, error) {
	path := fmt.Sprintf("/players/matchHistory/%s?includeDetails=true&limit=%d&offset=%d", playerID, limit, offset)
	var matches []Match
	if err := c.doRequest(ctx, path, &matches); err != nil {
		return nil, fmt.Errorf("get match history: %w", err)
	}
	return matches, nil
}

type Match struct {
	ID          string `json:"_id"`
	Version     string `json:"version"`
	QueueType   string `json:"queueType"`
	EndingWave  int    `json:"endingWave"`
	GameLength  int    `json:"gameLength"`
}
