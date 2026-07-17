package httpserver

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/yourname/legiontd2-copilot/internal/matrix"
	"github.com/yourname/legiontd2-copilot/internal/unitdata"
	"github.com/yourname/legiontd2-copilot/internal/ws"
)

func supplementUnitCosts(u *ws.HandUnit) {
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

func supplementState(s *ws.GameState) {
	for i := range s.Hand {
		supplementUnitCosts(&s.Hand[i])
	}
	for i := range s.Mercenaries {
		supplementUnitCosts(&s.Mercenaries[i])
	}
	for i := range s.TownActions {
		supplementUnitCosts(&s.TownActions[i])
	}
}

//go:embed static/*
var staticFS embed.FS

func New(addr string, hub *ws.Hub, iconsDir string) *http.Server {
	mux := http.NewServeMux()

	if iconsDir != "" {
		mux.Handle("/icons/", http.StripPrefix("/icons/", http.FileServer(http.Dir(iconsDir))))
	}

	mux.HandleFunc("/api/state", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		state := hub.GetState()
		supplementState(&state)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"state":           state,
			"recommendations": hub.GetRecs(),
			"matrix":          matrix.Build(state),
		})
	})

	mux.HandleFunc("/ws", hub.ServeWS)

	mux.HandleFunc("/api/ingest", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST required", 400)
			return
		}
		var state ws.GameState
		if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
			http.Error(w, "bad json", 400)
			return
		}
		hub.SetState(state)
		w.WriteHeader(204)
	})

	sub, err := fs.Sub(staticFS, "static")
	if err != nil {
		slog.Error("static fs sub", "error", err)
		panic(err)
	}
	mux.Handle("/", http.FileServer(http.FS(sub)))

	srv := &http.Server{Addr: addr, Handler: mux}
	go func() {
		slog.Info("HTTP+WS server", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("http server", "error", err)
		}
	}()
	return srv
}
