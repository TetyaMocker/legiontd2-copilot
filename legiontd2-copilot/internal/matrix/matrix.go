package matrix

import (
	"sort"

	"github.com/yourname/legiontd2-copilot/internal/unitdata"
	"github.com/yourname/legiontd2-copilot/internal/wavedata"
	"github.com/yourname/legiontd2-copilot/internal/ws"
)

type DamageCount struct {
	Normal int `json:"normal"`
	Pierce int `json:"pierce"`
	Magic  int `json:"magic"`
}

type RoleCount struct {
	Tank     int `json:"tank"`
	DPS      int `json:"dps"`
	Balanced int `json:"balanced"`
}

type EconomyFeatures struct {
	Gold        int     `json:"gold"`
	Mythium     int     `json:"mythium"`
	Supply      int     `json:"supply"`
	SupplyCap   int     `json:"supplyCap"`
	Income      int     `json:"income"`
	Wave        int     `json:"wave"`
	Phase       string  `json:"phase"`
	Timer       int     `json:"timer"`
	KingHP      float64 `json:"kingHp"`
	EnemyKingHP float64 `json:"enemyKingHp"`
	TeamGold    int     `json:"teamGold,omitempty"`
	EnemyTeamGold int   `json:"enemyTeamGold,omitempty"`
}

type BoardUnitInfo struct {
	Name       string  `json:"name"`
	Role       string  `json:"role"`
	DamageType string  `json:"damageType"`
	HP         float64 `json:"hp"`
}

type BoardAnalysis struct {
	TotalUnits int                `json:"totalUnits"`
	TotalHP    float64            `json:"totalHp"`
	ByDamage   DamageCount        `json:"byDamage"`
	ByRole     RoleCount          `json:"byRole"`
	Units      []BoardUnitInfo    `json:"units"`
}

type HandUnitInfo struct {
	Name       string `json:"name"`
	GoldCost   int    `json:"goldCost"`
	MythCost   int    `json:"mythCost"`
	SupplyCost int    `json:"supplyCost"`
	Role       string `json:"role"`
	DamageType string `json:"damageType"`
	Stacks     int    `json:"stacks"`
	Affordable bool   `json:"affordable"`
}

type AvailableAnalysis struct {
	Fighters []HandUnitInfo `json:"fighters"`
	Mercs    []HandUnitInfo `json:"mercs"`
}

type WaveForecast struct {
	Number      int     `json:"number"`
	Name        string  `json:"name"`
	ArmorType   string  `json:"armorType"`
	AttackType  string  `json:"attackType"`
	Amount      int     `json:"amount"`
	HasBoss     bool    `json:"hasBoss"`
	BestDamage  string  `json:"bestDamage"`
	Multiplier  float64 `json:"multiplier"`
}

type WaveAnalysis struct {
	Current  int            `json:"current"`
	Upcoming []WaveForecast `json:"upcoming"`
}

type CoverageAnalysis struct {
	BoardDamage     DamageCount `json:"boardDamage"`
	AvailableDamage DamageCount `json:"availableDamage"`
	MissingTypes    []string    `json:"missingTypes"`
	Recommended     string      `json:"recommended"`
	Explanation     string      `json:"explanation"`
}

type ContextFlags struct {
	IsFighting       bool `json:"isFighting"`
	IsKingLow        bool `json:"isKingLow"`
	HasAffordable    bool `json:"hasAffordable"`
	CanBuildWorkers  bool `json:"canBuildWorkers"`
}

type OpponentAnalysis struct {
	Players  []OpponentPlayer `json:"players"`
	Summary  OpponentSummary  `json:"summary"`
}

type OpponentPlayer struct {
	Name        string `json:"name"`
	Gold        int    `json:"gold"`
	Mythium     int    `json:"mythium"`
	Income      int    `json:"income"`
	Workers     int    `json:"workers"`
	Supply      int    `json:"supply"`
	SupplyCap   int    `json:"supplyCap"`
	FighterVal  int    `json:"fighterValue"`
	FieldUnits  int    `json:"fieldUnits"`
}

type OpponentSummary struct {
	TotalFieldUnits int        `json:"totalFieldUnits"`
	AvgGold         int        `json:"avgGold"`
	AvgMythium      int        `json:"avgMythium"`
	AvgIncome       int        `json:"avgIncome"`
	TotalFighterVal int        `json:"totalFighterValue"`
}

type FeatureMatrix struct {
	Economy   EconomyFeatures   `json:"economy"`
	Board     BoardAnalysis     `json:"board"`
	Available AvailableAnalysis `json:"available"`
	Waves     WaveAnalysis      `json:"waves"`
	Coverage  CoverageAnalysis  `json:"coverage"`
	Context   ContextFlags      `json:"context"`
	Opponent  OpponentAnalysis  `json:"opponent,omitempty"`
}

func Build(s ws.GameState) FeatureMatrix {
	m := FeatureMatrix{
		Economy: buildEconomy(s),
		Board:   buildBoard(s),
		Available: buildAvailable(s),
		Waves:   buildWaves(s),
		Context: buildContext(s),
	}
	m.Coverage = buildCoverage(s, m.Board, m.Available, m.Waves)
	m.Opponent = buildOpponent(s)
	return m
}

func buildEconomy(s ws.GameState) EconomyFeatures {
	return EconomyFeatures{
		Gold:        s.Gold,
		Mythium:     s.Mythium,
		Supply:      s.Supply,
		SupplyCap:   s.SupplyCap,
		Income:      s.Income,
		Wave:        s.Wave,
		Phase:       s.Phase,
		Timer:       s.WaveTimer,
		KingHP:      s.KingHP,
		EnemyKingHP: s.EnemyKingHP,
		TeamGold:    s.TeamGoldLeft,
		EnemyTeamGold: s.TeamGoldRight,
	}
}

func unitDamage(name string, isMerc bool) string {
	var at unitdata.AttackType
	var ok bool
	if isMerc {
		at, ok = unitdata.GetMercAttack(name)
	} else {
		at, ok = unitdata.GetFighterAttack(name)
	}
	if !ok {
		return "unknown"
	}
	return at.String()
}

func buildBoard(s ws.GameState) BoardAnalysis {
	b := BoardAnalysis{}
	for _, fu := range s.FieldUnits {
		b.TotalUnits++
		b.TotalHP += fu.HP
	}
	// Field units are anonymous (only ID + HP), so we can't determine
	// name/role/damageType without tracking build → unitId mapping.
	// Board analysis shows only counts and total HP.
	return b
}

func buildAvailable(s ws.GameState) AvailableAnalysis {
	a := AvailableAnalysis{}
	for _, u := range s.Hand {
		dt := unitDamage(u.Name, false)
		aff := u.Stacks > 0 && s.Gold >= u.CostGold && u.CostGold > 0 && (s.SupplyCap-s.Supply) >= u.CostSupply
		a.Fighters = append(a.Fighters, HandUnitInfo{
			Name:       u.Name,
			GoldCost:   u.CostGold,
			MythCost:   u.CostMythium,
			SupplyCost: u.CostSupply,
			Role:       u.Role,
			DamageType: dt,
			Stacks:     u.Stacks,
			Affordable: aff,
		})
	}
	for _, m := range s.Mercenaries {
		dt := unitDamage(m.Name, true)
		aff := m.Stacks > 0 && s.Mythium >= m.CostMythium && m.CostMythium > 0 && (s.SupplyCap-s.Supply) >= m.CostSupply
		a.Mercs = append(a.Mercs, HandUnitInfo{
			Name:       m.Name,
			GoldCost:   m.CostGold,
			MythCost:   m.CostMythium,
			SupplyCost: m.CostSupply,
			Role:       m.Role,
			DamageType: dt,
			Stacks:     m.Stacks,
			Affordable: aff,
		})
	}
	return a
}

func buildWaves(s ws.GameState) WaveAnalysis {
	w := WaveAnalysis{Current: s.Wave}
	start := s.Wave
	if s.Phase == "building" {
		start = s.Wave + 1
	}
	// Show up to 10 upcoming waves (starting from next wave)
	for n := start; n <= start+9 && n <= 21; n++ {
		wave, ok := wavedata.GetWave(n)
		if !ok {
			break
		}
		armor := unitdata.ParseArmor(string(wave.ArmorType))
		bestAtk := unitdata.BestAttackAgainst(armor)
		mult := unitdata.DamageMultiplier(bestAtk, armor)
		f := WaveForecast{
			Number:     n,
			Name:       wave.Name,
			ArmorType:  string(wave.ArmorType),
			AttackType: string(wave.AttackType),
			Amount:     wave.Amount + wave.Amount2,
			HasBoss:    wave.BossName != "",
			BestDamage: bestAtk.String(),
			Multiplier: mult,
		}
		w.Upcoming = append(w.Upcoming, f)
	}
	return w
}

func buildCoverage(s ws.GameState, board BoardAnalysis, avail AvailableAnalysis, waves WaveAnalysis) CoverageAnalysis {
	c := CoverageAnalysis{}

	// Count only AFFORDABLE available units by damage type
	for _, f := range avail.Fighters {
		if !f.Affordable {
			continue
		}
		switch f.DamageType {
		case "Normal":
			c.AvailableDamage.Normal++
		case "Pierce":
			c.AvailableDamage.Pierce++
		case "Magic":
			c.AvailableDamage.Magic++
		}
	}
	for _, m := range avail.Mercs {
		if !m.Affordable {
			continue
		}
		switch m.DamageType {
		case "Normal":
			c.AvailableDamage.Normal++
		case "Pierce":
			c.AvailableDamage.Pierce++
		case "Magic":
			c.AvailableDamage.Magic++
		}
	}
	// BoardDamage shows what's on field (anonymous, no damage types known)
	// So coverage is based on what can be built affordably right now
	c.BoardDamage = c.AvailableDamage

	// Determine what damage types are missing for upcoming waves
	bestDamageCount := map[string]int{}
	for _, wf := range waves.Upcoming {
		bestDamageCount[wf.BestDamage]++
	}
	// Sort by frequency
	type dc struct{ name string; count int }
	var sorted []dc
	for k, v := range bestDamageCount {
		sorted = append(sorted, dc{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].count > sorted[j].count })

	if len(sorted) > 0 {
		mostNeeded := sorted[0]
		// Use only available units (what can be built right now)
		availCount := 0
		switch mostNeeded.name {
		case "Normal":
			availCount = c.AvailableDamage.Normal
		case "Pierce":
			availCount = c.AvailableDamage.Pierce
		case "Magic":
			availCount = c.AvailableDamage.Magic
		}

		// Check gaps for each upcoming wave
		for _, wf := range waves.Upcoming {
			haveCount := 0
			switch wf.BestDamage {
			case "Normal":
				haveCount = c.AvailableDamage.Normal
			case "Pierce":
				haveCount = c.AvailableDamage.Pierce
			case "Magic":
				haveCount = c.AvailableDamage.Magic
			}
			if haveCount == 0 {
				c.MissingTypes = append(c.MissingTypes, wf.BestDamage)
			}
		}
		// Deduplicate
		seen := map[string]bool{}
		uniq := []string{}
		for _, t := range c.MissingTypes {
			if !seen[t] {
				seen[t] = true
				uniq = append(uniq, t)
			}
		}
		c.MissingTypes = uniq

		// Recommend
		if mostNeeded.count >= 2 && availCount < 2 {
			c.Recommended = mostNeeded.name
			c.Explanation = formatRecommendation(mostNeeded.name, mostNeeded.count, availCount)
		} else {
			c.Recommended = "balanced"
			c.Explanation = "Покрытие урона сбалансировано"
		}
		if len(c.MissingTypes) > 0 {
			if c.Recommended == "balanced" {
				c.Recommended = c.MissingTypes[0]
			}
			c.Explanation = "Не хватает урона типов: " + joinStrings(c.MissingTypes)
		}
	} else {
		c.Recommended = "unknown"
		c.Explanation = "Нет данных о волнах"
	}

	return c
}

func buildContext(s ws.GameState) ContextFlags {
	affordable := false
	for _, u := range s.Hand {
		if u.Stacks > 0 && s.Gold >= u.CostGold && u.CostGold > 0 && (s.SupplyCap-s.Supply) >= u.CostSupply {
			affordable = true
			break
		}
	}
	if !affordable {
		for _, m := range s.Mercenaries {
			if m.Stacks > 0 && s.Mythium >= m.CostMythium && m.CostMythium > 0 && (s.SupplyCap-s.Supply) >= m.CostSupply {
				affordable = true
				break
			}
		}
	}
	return ContextFlags{
		IsFighting:      s.Phase == "fighting",
		IsKingLow:       s.KingHP < 30.0,
		HasAffordable:   affordable,
		CanBuildWorkers: s.Gold >= 50 && s.Phase == "building",
	}
}

func formatRecommendation(dmg string, needCount, haveCount int) string {
	switch dmg {
	case "Normal":
		return "Броня врагов чаще всего уязвима к Normal урону. Построй больше Normal-юнитов"
	case "Pierce":
		return "Броня врагов чаще всего уязвима к Pierce урону. Построй больше Pierce-юнитов"
	case "Magic":
		return "Броня врагов чаще всего уязвима к Magic урону. Построй больше Magic-юнитов"
	}
	return ""
}

func joinStrings(s []string) string {
	r := ""
	for i, v := range s {
		if i > 0 {
			r += ", "
		}
		r += v
	}
	return r
}

func buildOpponent(s ws.GameState) OpponentAnalysis {
	oa := OpponentAnalysis{}
	if len(s.ScoreboardInfo) == 0 {
		return oa
	}
	// In 2v2: players 0&1 = our team (left), players 2&3 = enemy team (right)
	// ScoreboardInfo array index usually follows player number order
	for i, p := range s.ScoreboardInfo {
		pm, ok := p.(map[string]interface{})
		if !ok {
			continue
		}
		// Players 2+ are opponents (right team)
		if i < 2 {
			continue
		}
		pl := OpponentPlayer{
			Name: getStr(pm, "name"),
			Gold: getInt(pm, "gold"),
			Mythium: getInt(pm, "mythium"),
			Income: getInt(pm, "income"),
			Workers: getInt(pm, "workers"),
			Supply: getInt(pm, "supply"),
			SupplyCap: getInt(pm, "supplyCap"),
			FighterVal: getInt(pm, "value"),
		}
		// Count field units from grid
		if g, ok := pm["grid"].([]interface{}); ok {
			pl.FieldUnits = len(g)
		}
		oa.Players = append(oa.Players, pl)
	}
	// Build summary
	if len(oa.Players) > 0 {
		for _, pl := range oa.Players {
			oa.Summary.TotalFieldUnits += pl.FieldUnits
			oa.Summary.AvgGold += pl.Gold
			oa.Summary.AvgMythium += pl.Mythium
			oa.Summary.AvgIncome += pl.Income
			oa.Summary.TotalFighterVal += pl.FighterVal
		}
		n := len(oa.Players)
		oa.Summary.AvgGold /= n
		oa.Summary.AvgMythium /= n
		oa.Summary.AvgIncome /= n
	}
	return oa
}

func getStr(m map[string]interface{}, key string) string {
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}

func getInt(m map[string]interface{}, key string) int {
	v, ok := m[key]
	if !ok || v == nil {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	}
	return 0
}
