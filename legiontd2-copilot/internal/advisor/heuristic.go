package advisor

type EconomySnapshot struct {
	Mythium           int
	Income            int
	WaveNumber        int
	WaveTimerSeconds  int
	KingHPPercent     int
	AllyKingHPPercent int
	Confidence        float32
}

type Recommendation struct {
	Kind        string
	Explanation string
}

type Advisor interface {
	Recommend(EconomySnapshot) []Recommendation
}

type HeuristicAdvisor struct{}

func NewHeuristicAdvisor() *HeuristicAdvisor {
	return &HeuristicAdvisor{}
}

func (h *HeuristicAdvisor) Recommend(snap EconomySnapshot) []Recommendation {
	var recs []Recommendation

	if snap.Confidence < 0.3 {
		recs = append(recs, Recommendation{
			Kind:        "save",
			Explanation: "Распознавание ненадёжно — рекомендации временно отключены",
		})
		return recs
	}

	if snap.Mythium > 120 && snap.WaveTimerSeconds <= 10 {
		recs = append(recs, Recommendation{
			Kind:        "spend",
			Explanation: "До волны осталось мало времени — потрать Mythium на юнитов/апгрейды",
		})
	}

	if snap.Mythium > 200 && snap.WaveTimerSeconds > 15 {
		recs = append(recs, Recommendation{
			Kind:        "save",
			Explanation: "Рано тратить — копи Mythium для отправки наёмников на волне противника",
		})
	}

	if snap.KingHPPercent < 30 && snap.Mythium > 60 {
		recs = append(recs, Recommendation{
			Kind:        "spend",
			Explanation: "HP короля низкое — инвестируй в защиту",
		})
	}

	if len(recs) == 0 {
		recs = append(recs, Recommendation{
			Kind:        "save",
			Explanation: "Ситуация стабильная — продолжай копить",
		})
	}

	return recs
}
