package ws

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type GameState struct {
	Mythium              int     `json:"mythium"`
	Gold                 int     `json:"gold"`
	Supply               int     `json:"supply"`
	SupplyCap            int     `json:"supplyCap"`
	Income               int     `json:"income"`
	MythiumGatherRate    float64 `json:"mythiumGatherRate"`
	EstimatedMythium     int     `json:"estimatedMythium"`
	Wave                 int     `json:"wave"`
	WaveTimer            int     `json:"waveTimer"`
	KingHP               float64 `json:"kingHp"`
	EnemyKingHP          float64 `json:"enemyKingHp"`
	EnemiesRemainingWest int     `json:"enemiesRemainingWest"`
	EnemiesRemainingEast int     `json:"enemiesRemainingEast"`
	TimeElapsed          float64 `json:"timeElapsed"`
	Phase                string  `json:"phase"`
}

type Recommendation struct {
	Action   string `json:"action"`
	Message  string `json:"message"`
	Priority int    `json:"priority"`
}

type Hub struct {
	mu        sync.RWMutex
	state     GameState
	recs      []Recommendation
	clients   map[*websocket.Conn]bool
}

func NewHub() *Hub {
	return &Hub{clients: make(map[*websocket.Conn]bool)}
}

func (h *Hub) GetState() GameState {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.state
}

func (h *Hub) SetState(s GameState) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.state = s
}

func (h *Hub) GetRecs() []Recommendation {
	h.mu.RLock()
	defer h.mu.RUnlock()
	recs := make([]Recommendation, len(h.recs))
	copy(recs, h.recs)
	return recs
}

func (h *Hub) SetRecs(recs []Recommendation) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.recs = recs
	h.broadcast(recs)
}

func (h *Hub) broadcast(recs []Recommendation) {
	msg, err := json.Marshal(map[string]interface{}{
		"type": "recommendation",
		"data": recs,
	})
	if err != nil {
		return
	}
	for conn := range h.clients {
		err := conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			slog.Warn("ws write", "error", err)
			conn.Close()
			delete(h.clients, conn)
		}
	}
}

func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("ws upgrade", "error", err)
		return
	}

	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()

	// Send current recs immediately
	if recs := h.GetRecs(); len(recs) > 0 {
		msg, _ := json.Marshal(map[string]interface{}{
			"type": "recommendation",
			"data": recs,
		})
		conn.WriteMessage(websocket.TextMessage, msg)
	}

	pingTicker := time.NewTicker(30 * time.Second)
	defer func() {
		pingTicker.Stop()
		h.mu.Lock()
		delete(h.clients, conn)
		h.mu.Unlock()
		conn.Close()
	}()

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		select {
		case <-pingTicker.C:
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		default:
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			var state GameState
			if err := json.Unmarshal(msg, &state); err == nil {
				h.SetState(state)
			}
		}
	}
}
