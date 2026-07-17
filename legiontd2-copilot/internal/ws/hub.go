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
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

type HandUnit struct {
	ActionID    int    `json:"actionId"`
	Name       string `json:"name"`
	Icon       string `json:"icon"`
	CostGold   int    `json:"costGold"`
	CostMythium int   `json:"costMythium"`
	CostSupply int    `json:"costSupply"`
	Stacks     int    `json:"stacks"`
	Role       string `json:"role"`
}

type FieldUnit struct {
	UnitID    int     `json:"unitId"`
	HP        float64 `json:"hp"`
	FirstSeen int64   `json:"firstSeen"`
}

type BuildQueueItem struct {
	ActionID int `json:"actionId"`
	Count    int `json:"count"`
}

type PlayerScore struct {
	Player int             `json:"player"`
	Value  json.RawMessage `json:"value"`
}

type SelectedUnit struct {
	Name  string  `json:"name"`
	HP    float64 `json:"hp"`
	MaxHP float64 `json:"maxHp"`
	Title string  `json:"title"`
}

type Recommendation struct {
	Action   string `json:"action"`
	Message  string `json:"message"`
	Priority int    `json:"priority"`
}

type GameState struct {
	Mythium              int              `json:"mythium"`
	Gold                 int              `json:"gold"`
	Supply               int              `json:"supply"`
	SupplyCap            int              `json:"supplyCap"`
	Income               int              `json:"income"`
	MythiumGatherRate    float64          `json:"mythiumGatherRate"`
	EstimatedMythium     int              `json:"estimatedMythium"`
	Wave                 int              `json:"wave"`
	WaveTimer            int              `json:"waveTimer"`
	KingHP               float64          `json:"kingHp"`
	EnemyKingHP          float64          `json:"enemyKingHp"`
	EnemiesRemainingWest int              `json:"enemiesRemainingWest"`
	EnemiesRemainingEast int              `json:"enemiesRemainingEast"`
	TimeElapsed          float64          `json:"timeElapsed"`
	Phase                string           `json:"phase"`

	Hand         []HandUnit        `json:"hand"`
	FieldUnits   map[int]FieldUnit `json:"fieldUnits"`
	BuildQueue   []BuildQueueItem  `json:"buildQueue"`
	PurchaseQueue []BuildQueueItem `json:"purchaseQueue"`
	Mercenaries  []HandUnit        `json:"mercenaries"`
	TownActions  []HandUnit        `json:"townActions"`
	PlayerScores []PlayerScore     `json:"playerScores"`
	SelectedUnit *SelectedUnit     `json:"selectedUnit,omitempty"`

	ActionSample     interface{} `json:"_actionSample,omitempty"`
	MercActionSample interface{} `json:"_mercActionSample,omitempty"`

	ScoreboardInfo  []interface{} `json:"scoreboardInfo,omitempty"`
	TeamGoldLeft    int           `json:"teamGoldLeft,omitempty"`
	TeamGoldRight   int           `json:"teamGoldRight,omitempty"`
	KingUpgradesLeft  interface{} `json:"kingUpgradesLeft,omitempty"`
	KingUpgradesRight interface{} `json:"kingUpgradesRight,omitempty"`
	Moneylender       interface{} `json:"moneylender,omitempty"`
}

type Hub struct {
	mu        sync.RWMutex
	state     GameState
	recs      []Recommendation
	clients   map[*websocket.Conn]bool
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]bool),
		state: GameState{
			FieldUnits:   make(map[int]FieldUnit),
			Hand:         make([]HandUnit, 0),
			BuildQueue:   make([]BuildQueueItem, 0),
			Mercenaries:  make([]HandUnit, 0),
			TownActions:  make([]HandUnit, 0),
			PlayerScores: make([]PlayerScore, 0),
		},
	}
}

func (h *Hub) GetState() GameState {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.state
}

func (h *Hub) SetState(s GameState) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if s.FieldUnits == nil {
		s.FieldUnits = make(map[int]FieldUnit)
	}
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
			} else {
				slog.Warn("ws unmarshal", "error", err)
			}
		}
	}
}
