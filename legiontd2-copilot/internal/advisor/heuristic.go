package advisor

import (
	"fmt"
	"strings"

	"github.com/yourname/legiontd2-copilot/internal/ws"
)

func Recommend(state ws.GameState) []ws.Recommendation {
	var recs []ws.Recommendation

	if state.Wave == 0 && state.Mythium == 0 && len(state.Hand) == 0 {
		return recs
	}

	affordableUnits := make([]ws.HandUnit, 0)
	for _, u := range state.Hand {
		if u.Stacks > 0 && state.Gold >= u.CostGold && state.SupplyCap-state.Supply >= u.CostSupply {
			affordableUnits = append(affordableUnits, u)
		}
	}

	if len(state.Hand) > 0 && len(affordableUnits) == 0 && state.WaveTimer > 0 {
		recs = append(recs, ws.Recommendation{
			Action:   "save_gold",
			Message:  "Не хватает золота/снабжения для постановки юнитов — копи",
			Priority: 1,
		})
	}

	if len(affordableUnits) > 0 && state.WaveTimer <= 30 {
		var names []string
		for _, u := range affordableUnits {
			names = append(names, u.Name)
		}
		recs = append(recs, ws.Recommendation{
			Action:   "build_units",
			Message:  fmt.Sprintf("Поставь доступные юниты: %s", strings.Join(names, ", ")),
			Priority: 3,
		})
	}

	if state.WaveTimer <= 10 && state.Mythium > 80 {
		recs = append(recs, ws.Recommendation{
			Action:   "spend_mythium",
			Message:  "До волны осталось мало времени — отправь наёмников",
			Priority: 4,
		})
	}

	if state.Mythium > 120 && state.WaveTimer > 15 {
		recs = append(recs, ws.Recommendation{
			Action:   "save_mythium",
			Message:  "Копи Mythium для отправки наёмников на волне противника",
			Priority: 2,
		})
	}

	if state.KingHP < 30 && state.Mythium > 60 {
		recs = append(recs, ws.Recommendation{
			Action:   "upgrade_king",
			Message:  "HP короля низкое — улучши короля для защиты",
			Priority: 5,
		})
	}

	if len(recs) == 0 {
		recs = append(recs, ws.Recommendation{
			Action:   "hold",
			Message:  "Ситуация стабильная — продолжай развиваться",
			Priority: 1,
		})
	}

	return recs
}
