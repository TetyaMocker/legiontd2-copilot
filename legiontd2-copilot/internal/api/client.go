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
		httpClient: &http.Client{Timeout: 15 * time.Second},
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

func (c *Client) GetUnitByName(ctx context.Context, name string, version string) (*UnitStats, error) {
	path := "/units/byName/" + name
	if version != "" {
		path += "?version=" + version
	}
	var u UnitStats
	if err := c.doRequest(ctx, path, &u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (c *Client) GetUnitsByVersion(ctx context.Context, version string, limit, offset int) ([]UnitStats, error) {
	path := fmt.Sprintf("/units/byVersion/%s?limit=%d&offset=%d", version, limit, offset)
	var units []UnitStats
	if err := c.doRequest(ctx, path, &units); err != nil {
		return nil, err
	}
	return units, nil
}

func (c *Client) GetPlayerByName(ctx context.Context, name string) (*Player, error) {
	var players []Player
	if err := c.doRequest(ctx, "/players/byName/"+name, &players); err != nil {
		return nil, err
	}
	if len(players) == 0 {
		return nil, fmt.Errorf("player %q not found", name)
	}
	return &players[0], nil
}

func (c *Client) GetPlayerByID(ctx context.Context, id string) (*Player, error) {
	var players []Player
	if err := c.doRequest(ctx, "/players/byId/"+id, &players); err != nil {
		return nil, err
	}
	if len(players) == 0 {
		return nil, fmt.Errorf("player %q not found", id)
	}
	return &players[0], nil
}

func (c *Client) GetPlayerStats(ctx context.Context, id string) (*Stats, error) {
	var s Stats
	if err := c.doRequest(ctx, "/players/stats/"+id, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (c *Client) GetTopPlayers(ctx context.Context, sortBy string, limit, offset int) ([]Stats, error) {
	if sortBy == "" {
		sortBy = "overallElo"
	}
	path := fmt.Sprintf("/players/stats?sortBy=%s&limit=%d&offset=%d&sortDirection=-1", sortBy, limit, offset)
	var stats []Stats
	if err := c.doRequest(ctx, path, &stats); err != nil {
		return nil, err
	}
	return stats, nil
}

func (c *Client) GetMatchHistory(ctx context.Context, playerID string, limit, offset int, includeDetails bool) ([]Match, error) {
	path := fmt.Sprintf("/players/matchHistory/%s?limit=%d&offset=%d", playerID, limit, offset)
	if includeDetails {
		path += "&includeDetails=true"
	}
	var matches []Match
	if err := c.doRequest(ctx, path, &matches); err != nil {
		return nil, err
	}
	return matches, nil
}

func (c *Client) GetMatchByID(ctx context.Context, id string, includeDetails bool) (*Match, error) {
	path := "/games/byId/" + id
	if includeDetails {
		path += "?includeDetails=true"
	}
	var m Match
	if err := c.doRequest(ctx, path, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (c *Client) GetMatchesByFilter(ctx context.Context, version string, limit, offset int, includeDetails bool) ([]Match, error) {
	path := fmt.Sprintf("/games?limit=%d&offset=%d", limit, offset)
	if version != "" {
		path += "&version=" + version
	}
	if includeDetails {
		path += "&includeDetails=true"
	}
	var matches []Match
	if err := c.doRequest(ctx, path, &matches); err != nil {
		return nil, err
	}
	return matches, nil
}

func (c *Client) GetWaves(ctx context.Context, limit, offset int) ([]Wave, error) {
	path := fmt.Sprintf("/info/waves/%d/%d", offset, limit)
	var waves []Wave
	if err := c.doRequest(ctx, path, &waves); err != nil {
		return nil, err
	}
	return waves, nil
}

func (c *Client) GetLegions(ctx context.Context, limit, offset int, playable bool) ([]Legion, error) {
	path := fmt.Sprintf("/info/legions/%d/%d?playable=%t", offset, limit, playable)
	var legions []Legion
	if err := c.doRequest(ctx, path, &legions); err != nil {
		return nil, err
	}
	return legions, nil
}
