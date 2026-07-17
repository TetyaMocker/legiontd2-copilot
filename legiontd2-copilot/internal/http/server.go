package httpserver

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/yourname/legiontd2-copilot/internal/ws"
)

//go:embed static/*
var staticFS embed.FS

func New(addr string, hub *ws.Hub) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/state", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"state":           hub.GetState(),
			"recommendations": hub.GetRecs(),
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
