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

type IngestPayload struct {
	Mythium    int     `json:"mythium"`
	Income     int     `json:"income"`
	Wave       int     `json:"wave"`
	Timer      int     `json:"timer"`
	KingHP     int     `json:"king_hp"`
	Confidence float64 `json:"confidence"`
}

type StateJSON struct {
	Mythium     int       `json:"mythium"`
	Income      int       `json:"income"`
	Wave        int       `json:"wave"`
	WaveTimer   int       `json:"waveTimer"`
	KingHP      int       `json:"kingHp"`
	Confidence  float64   `json:"confidence"`
	Recs        []recJSON `json:"recommendations"`
}

type recJSON struct {
	Kind        string `json:"kind"`
	Explanation string `json:"explanation"`
}

type AppState struct {
	mu   sync.RWMutex
	eco  IngestPayload
	recs []advisor.Recommendation
}

func (a *AppState) Ingest(payload IngestPayload) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.eco = payload
}

func (a *AppState) SetRecs(recs []advisor.Recommendation) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.recs = recs
}

func (a *AppState) Snapshot() StateJSON {
	a.mu.RLock()
	defer a.mu.RUnlock()
	recs := make([]recJSON, len(a.recs))
	for i, r := range a.recs {
		recs[i] = recJSON{Kind: r.Kind, Explanation: r.Explanation}
	}
	return StateJSON{
		Mythium:    a.eco.Mythium,
		Income:     a.eco.Income,
		Wave:       a.eco.Wave,
		WaveTimer:  a.eco.Timer,
		KingHP:     a.eco.KingHP,
		Confidence: a.eco.Confidence,
		Recs:       recs,
	}
}

func New(state *AppState) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/state", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(state.Snapshot())
	})

	mux.HandleFunc("/api/ingest", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST required", 400)
			return
		}
		var payload IngestPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "bad json: "+err.Error(), 400)
			return
		}
		state.Ingest(payload)
		w.WriteHeader(204)
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
		slog.Info("web UI + ingestion server", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
		}
	}()
	return srv
}
