package webserver

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log/slog"
	"net/http"
	"sync"

	"github.com/yourname/legiontd2-copilot/internal/advisor"
)

//go:embed static/*
var staticFS embed.FS

type AppState struct {
	mu          sync.RWMutex
	Economy     advisor.EconomySnapshot
	Recs        []advisor.Recommendation
}

type StateJSON struct {
	Mythium     int     `json:"mythium"`
	Income      int     `json:"income"`
	Wave        int     `json:"wave"`
	WaveTimer   int     `json:"waveTimer"`
	KingHP      int     `json:"kingHp"`
	Confidence  float32 `json:"confidence"`
	Recs        []recJSON `json:"recommendations"`
}

type recJSON struct {
	Kind        string `json:"kind"`
	Explanation string `json:"explanation"`
}

func (a *AppState) Update(eco advisor.EconomySnapshot, recs []advisor.Recommendation) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Economy = eco
	a.Recs = recs
}

func (a *AppState) Snapshot() StateJSON {
	a.mu.RLock()
	defer a.mu.RUnlock()
	recs := make([]recJSON, len(a.Recs))
	for i, r := range a.Recs {
		recs[i] = recJSON{Kind: r.Kind, Explanation: r.Explanation}
	}
	return StateJSON{
		Mythium:    a.Economy.Mythium,
		Income:     a.Economy.Income,
		Wave:       a.Economy.WaveNumber,
		WaveTimer:  a.Economy.WaveTimerSeconds,
		KingHP:     a.Economy.KingHPPercent,
		Confidence: a.Economy.Confidence,
		Recs:       recs,
	}
}

func New(state *AppState) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/state", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(state.Snapshot())
	})

	sub, err := fs.Sub(staticFS, "static")
	if err != nil {
		slog.Error("static fs sub", "error", err)
		panic(err)
	}
	mux.Handle("/", http.FileServer(http.FS(sub)))

	return mux
}

func Start(state *AppState, addr string) *http.Server {
	srv := &http.Server{
		Addr:    addr,
		Handler: New(state),
	}
	go func() {
		slog.Info("web UI server started", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("web server error", "error", err)
		}
	}()
	return srv
}
