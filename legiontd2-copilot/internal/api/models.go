package api

type Player struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}

type Stats struct {
	PlayerID  string `json:"playerId"`
	PlayerName string `json:"playerName"`
	OverallElo int    `json:"overallElo"`
	ClassicElo int    `json:"classicElo"`
	Wins       int    `json:"wins"`
	Losses     int    `json:"losses"`
	GamesPlayed int   `json:"gamesPlayed"`
}

type Match struct {
	ID          string              `json:"_id"`
	Version     string              `json:"version"`
	Date        string              `json:"date"`
	QueueType   string              `json:"queueType"`
	EndingWave  int                 `json:"endingWave"`
	GameLength  int                 `json:"gameLength"`
	GameElo     int                 `json:"gameElo"`
	PlayerCount int                 `json:"playerCount"`
	HumanCount  int                 `json:"humanCount"`
	PlayersData []PlayerMatchDetails `json:"playersData,omitempty"`
}

type PlayerMatchDetails struct {
	PlayerID    string   `json:"playerId"`
	PlayerName  string   `json:"playerName"`
	PlayerSlot  int      `json:"playerSlot"`
	Legion      string   `json:"legion"`
	Workers     int      `json:"workers"`
	Value       int      `json:"value"`
	GameResult  string   `json:"gameResult"`
	OverallElo  int      `json:"overallElo"`
	ClassicElo  int      `json:"classicElo"`
	Fighters    string   `json:"fighters"`
	Mercenaries string   `json:"mercenaries"`
	FirstWaveFighters string `json:"firstWaveFighters"`
	Rolls       string   `json:"rolls"`

	NetWorthPerWave  []int    `json:"netWorthPerWave"`
	WorkersPerWave   []int    `json:"workersPerWave"`
	IncomePerWave    []int    `json:"incomePerWave"`
	BuildPerWave     []string `json:"buildPerWave"`
	LeaksPerWave     []string `json:"leaksPerWave"`
	MercenariesSentPerWave  []any  `json:"mercenariesSentPerWave"`
	MercenariesReceivedPerWave []any `json:"mercenariesReceivedPerWave"`
	KingUpgradesPerWave  []string `json:"kingUpgradesPerWave"`
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
