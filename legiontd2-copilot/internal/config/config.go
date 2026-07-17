package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Region struct {
	X     int   `json:"x"`
	Y     int   `json:"y"`
	W     int   `json:"w"`
	H     int   `json:"h"`
	Label string `json:"label"`
	Range []int `json:"range,omitempty"`
}

type CaptureConfig struct {
	Mythium Region `json:"mythium"`
	Income  Region `json:"income"`
	Wave    Region `json:"wave"`
	Timer   Region `json:"timer"`
	KingHP  Region `json:"king_hp"`
}

type TrackingConfig struct {
	ReScanIntervalMs int     `json:"re_scan_interval_ms"`
	ChangeThreshold  float64 `json:"change_threshold"`
	MaxCacheAgeMs    int     `json:"max_cache_age_ms"`
}

type Config struct {
	SchemaVersion      int            `json:"schema_version"`
	PatchVersion       string         `json:"patch_version"`
	Resolution         string         `json:"resolution"`
	ConfidenceThreshold float64       `json:"confidence_threshold"`
	OCRTimeoutMs       int            `json:"ocr_timeout_ms"`
	Capture            CaptureConfig  `json:"capture"`
	Tracking           TrackingConfig `json:"tracking"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}

func LoadDefault() *Config {
	return &Config{
		SchemaVersion: 1,
		PatchVersion:  "any",
		Resolution:    "1920x1080",
		ConfidenceThreshold: 0.7,
		OCRTimeoutMs:  5000,
		Capture: CaptureConfig{
			Mythium: Region{X: 520, Y: 50, W: 120, H: 35, Label: "Mythium", Range: []int{0, 99999}},
			Income:  Region{X: 515, Y: 108, W: 80, H: 28, Label: "Income", Range: []int{0, 9999}},
			Wave:    Region{X: 850, Y: 90, W: 160, H: 35, Label: "Wave number", Range: []int{1, 30}},
			Timer:   Region{X: 890, Y: 58, W: 100, H: 30, Label: "Wave timer", Range: []int{0, 600}},
			KingHP:  Region{X: 1410, Y: 85, W: 70, H: 35, Label: "King HP", Range: []int{0, 100}},
		},
		Tracking: TrackingConfig{
			ReScanIntervalMs: 3000,
			ChangeThreshold:  0.15,
			MaxCacheAgeMs:    10000,
		},
	}
}
