package advisor

import (
	"fmt"
	"strings"

	"github.com/yourname/legiontd2-copilot/internal/unitdata"
	"github.com/yourname/legiontd2-copilot/internal/wavedata"
	"github.com/yourname/legiontd2-copilot/internal/ws"
)

func Recommend(state ws.GameState) []ws.Recommendation {
	var recs []ws.Recommendation

	if state.Wave == 0 && state.Mythium == 0 && len(state.Hand) == 0 {
		return recs
	}

	nextWaveNum := state.Wave
	if state.Phase == "building" {
		nextWaveNum = state.Wave + 1
	}

	nextWave, hasWave := wavedata.GetWave(nextWaveNum)

	affordable := make([]ws.HandUnit, 0)
	for _, u := range state.Hand {
		if u.Stacks > 0 && state.Gold >= u.CostGold && u.CostGold > 0 && state.SupplyCap-state.Supply >= u.CostSupply {
			affordable = append(affordable, u)
		}
	}

	canWork := state.Gold >= 50 && state.Phase == "building"
	lowHealth := state.KingHP < 30.0

	if hasWave && state.Phase == "building" && nextWaveNum <= 21 {
		enemyArmor := unitdata.ParseArmor(string(nextWave.ArmorType))
		bestAtk := unitdata.BestAttackAgainst(enemyArmor)

		counterUnits := make([]string, 0)
		for _, u := range state.Hand {
			if u.Stacks > 0 {
				uat, ok := unitdata.GetFighterAttack(u.Name)
				if ok && uat == bestAtk {
					counterUnits = append(counterUnits, u.Name)
				}
			}
		}

		if len(counterUnits) > 0 {
			recs = append(recs, ws.Recommendation{
				Action:   "counter_wave",
				Message:  fmt.Sprintf("Волна %d (%s) — броня %s, эффективен %s урон. Есть в руке: %s", nextWaveNum, nextWave.Name, string(nextWave.ArmorType), bestAtk, strings.Join(counterUnits, ", ")),
				Priority: 4,
			})
		} else {
			recs = append(recs, ws.Recommendation{
				Action:   "wave_info",
				Message:  fmt.Sprintf("Волна %d (%s) — броня %s, эффективен %s урон. Нет контр-юнитов в руке", nextWaveNum, nextWave.Name, string(nextWave.ArmorType), bestAtk),
				Priority: 2,
			})
		}

		if nextWave.ArmorType == wavedata.ArmorLight && bestAtk == unitdata.AtkPierce {
			multi := unitdata.DamageMultiplier(unitdata.AtkPierce, enemyArmor)
			recs = append(recs, ws.Recommendation{
				Action:   "tip",
				Message:  fmt.Sprintf("Light броня — Pierce урон на %.0f%% эффективнее", (multi-1)*100),
				Priority: 1,
			})
		}
		if nextWave.ArmorType == wavedata.ArmorFortified && bestAtk == unitdata.AtkNormal {
			recs = append(recs, ws.Recommendation{
				Action:   "tip",
				Message:  fmt.Sprintf("Fortified броня — Normal урон наиболее эффективен"),
				Priority: 1,
			})
		}
	}

	if state.WaveTimer > 0 && state.WaveTimer <= 30 && len(affordable) > 0 {
		var names []string
		for _, u := range affordable {
			names = append(names, u.Name)
		}
		recs = append(recs, ws.Recommendation{
			Action:   "build_units",
			Message:  fmt.Sprintf("До волны %ds — поставь: %s", state.WaveTimer, strings.Join(names, ", ")),
			Priority: 5,
		})
	}

	if state.WaveTimer > 10 && canWork && len(state.Hand) == 0 && state.SupplyCap-state.Supply < 5 {
		recs = append(recs, ws.Recommendation{
			Action:   "build_workers",
			Message:  "Нет юнитов для постановки — создай рабочих (+10 дохода)",
			Priority: 4,
		})
	} else if state.WaveTimer > 10 && canWork && state.Gold >= 100 && state.Wave <= 5 {
		recs = append(recs, ws.Recommendation{
			Action:   "build_workers",
			Message:  fmt.Sprintf("Ранняя игра — создай рабочих (есть %d золота)", state.Gold),
			Priority: 3,
		})
	}

	if state.WaveTimer > 0 && state.WaveTimer <= 10 && state.Mythium > 60 {
		recs = append(recs, ws.Recommendation{
			Action:   "send_mercs",
			Message:  fmt.Sprintf("Волна вот-вот начнётся — отправь наёмников (%d мифиума)", state.Mythium),
			Priority: 5,
		})
	} else if state.Mythium > 120 && state.Phase == "building" {
		recs = append(recs, ws.Recommendation{
			Action:   "save_mythium",
			Message:  fmt.Sprintf("Копи мифиум (%d) — отправь наёмников на волне противника", state.Mythium),
			Priority: 3,
		})
	}

	if lowHealth && state.Mythium > 40 && state.Phase == "building" {
		recs = append(recs, ws.Recommendation{
			Action:   "upgrade_king",
			Message:  fmt.Sprintf("HP короля %.0f%% — улучши короля для защиты", state.KingHP),
			Priority: 5,
		})
	}

	if state.Phase == "building" && len(state.Hand) > 0 && len(affordable) == 0 && state.Gold < 50 && state.Mythium < 60 {
		recs = append(recs, ws.Recommendation{
			Action:   "save_gold",
			Message:  fmt.Sprintf("Мало ресурсов — жди дохода (золото: %d, мифиум: %d)", state.Gold, state.Mythium),
			Priority: 2,
		})
	}

	if len(recs) == 0 {
		if state.Phase == "building" {
			recs = append(recs, ws.Recommendation{
				Action:   "hold",
				Message:  "Стройся — ситуация стабильная",
				Priority: 1,
			})
		} else {
			recs = append(recs, ws.Recommendation{
				Action:   "fight",
				Message:  "Идёт бой — наблюдает за волной",
				Priority: 1,
			})
		}
	}

	recs = append(recs, ws.Recommendation{
		Action:   "save_mythium",
		Message:  fmt.Sprintf("Mythium: %d | Доход: %d | Золото: %d | Снабжение: %d/%d | Волна: %d", state.Mythium, state.Income, state.Gold, state.Supply, state.SupplyCap, state.Wave),
		Priority: 0,
	})

	return recs
}
