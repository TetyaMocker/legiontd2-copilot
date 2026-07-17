package advisor

import (
	"github.com/yourname/legiontd2-copilot/internal/ws"
)

func Recommend(state ws.GameState) []ws.Recommendation {
	var recs []ws.Recommendation

	if state.Wave == 0 && state.Mythium == 0 {
		return recs
	}

	if state.WaveTimer <= 10 && state.Mythium > 80 {
		recs = append(recs, ws.Recommendation{
			Action:   "spend_mythium",
			Message:  "До волны осталось мало времени — потрать Mythium на наёмников",
			Priority: 3,
		})
	}

	if state.Mythium > 120 && state.WaveTimer > 15 {
		recs = append(recs, ws.Recommendation{
			Action:   "save_mythium",
			Message:  "Копи Mythium — отправь наёмников на волне противника",
			Priority: 2,
		})
	}

	if state.KingHP < 30 && state.Mythium > 60 {
		recs = append(recs, ws.Recommendation{
			Action:   "upgrade_king",
			Message:  "HP короля низкое — улучши короля",
			Priority: 5,
		})
	}

	if len(recs) == 0 {
		recs = append(recs, ws.Recommendation{
			Action:   "hold",
			Message:  "Ситуация стабильная — продолжай копить",
			Priority: 1,
		})
	}

	return recs
}
