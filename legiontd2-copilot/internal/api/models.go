package api

import "encoding/json"

type Player struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}

type Stats struct {
	PlayerID   string `json:"_id"`
	PlayerName string `json:"playerName,omitempty"`
	OverallElo int    `json:"overallElo"`
	ClassicElo int    `json:"classicElo"`
	Wins       int    `json:"wins"`
	Losses     int    `json:"losses"`
	GamesPlayed int   `json:"gamesPlayed"`
}

type Match struct {
	ID          string          `json:"_id"`
	Version     string          `json:"version"`
	Date        string          `json:"date"`
	QueueType   string          `json:"queueType"`
	EndingWave  int             `json:"endingWave"`
	GameLength  int             `json:"gameLength"`
	GameElo     int             `json:"gameElo"`
	PlayerCount int             `json:"playerCount"`
	HumanCount  int             `json:"humanCount"`
	PlayersData json.RawMessage `json:"playersData,omitempty"`
}

type PlayerMatchDetails struct {
	PlayerID    string   `json:"playerId"`
	PlayerName  string   `json:"playerName"`
	PlayerSlot  int      `json:"playerSlot"`
	Legion      string   `json:"legion"`
	Workers     float64  `json:"workers"`
	Value       int      `json:"value"`
	GameResult  string   `json:"gameResult"`
	OverallElo  int      `json:"overallElo"`
	ClassicElo  int      `json:"classicElo"`
	EloChange   int      `json:"eloChange"`
	Fighters    string   `json:"fighters"`
	Mercenaries string   `json:"mercenaries"`
	FirstWaveFighters string `json:"firstWaveFighters"`
	Rolls       string   `json:"rolls"`

	NetWorthPerWave  json.RawMessage `json:"netWorthPerWave"`
	ValuePerWave     json.RawMessage `json:"valuePerWave"`
	WorkersPerWave   json.RawMessage `json:"workersPerWave"`
	IncomePerWave    json.RawMessage `json:"incomePerWave"`
	MercenariesSentPerWave  json.RawMessage `json:"mercenariesSentPerWave"`
	MercenariesReceivedPerWave json.RawMessage `json:"mercenariesReceivedPerWave"`
	LeaksPerWave     json.RawMessage `json:"leaksPerWave"`
	BuildPerWave     json.RawMessage `json:"buildPerWave"`
	KingUpgradesPerWave     json.RawMessage `json:"kingUpgradesPerWave"`
	OpponentKingUpgradesPerWave json.RawMessage `json:"opponentKingUpgradesPerWave"`
	ChosenSpell      string `json:"chosenSpell"`
	MvpScore         int    `json:"mvpScore"`
	LeakValue        int    `json:"leakValue"`
	LeaksCaughtValue int    `json:"leaksCaughtValue"`
	PartySize        int    `json:"partySize"`
	StayedUntilEnd   bool   `json:"stayedUntilEnd"`
	Doubledown       bool   `json:"doubledown"`
}

type UnitStats struct {
	ID        string `json:"_id"`
	Name      string `json:"name"`
	HP        string `json:"hp"`
	GoldCost  string `json:"goldCost"`
	MythiumCost string `json:"mythiumCost"`
	ArmorType string `json:"armorType"`
	UnitClass string `json:"unitClass"`
	InfoTier  string `json:"infoTier"`
	LegionID  string `json:"legionId"`
	IconPath  string `json:"iconPath"`
}

type Wave struct {
	ID    string `json:"_id"`
	Name  string `json:"name"`
	Level int    `json:"levelnum"`
	Unit  string `json:"unit"`
	Amount int   `json:"amount"`
	TotalReward int `json:"totalReward"`
	PrepareTime int `json:"prepareTime"`
}

type Legion struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}

type GamesResponse struct {
	Value []Match `json:"value"`
	Count int     `json:"Count"`
}
