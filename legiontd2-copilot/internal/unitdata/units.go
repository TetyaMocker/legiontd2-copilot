package unitdata

type AttackType int

const (
	AtkNormal AttackType = iota
	AtkPierce
	AtkMagic
)

// UnitCost holds game resource costs for spawning a unit.
type UnitCost struct {
	Gold    int
	Mythium int
	Supply  int
}

var FighterCosts = map[string]UnitCost{}

var MercCosts = map[string]UnitCost{
	"Witch":        {Mythium: 30, Supply: 2},
	"Siege Ram":    {Mythium: 45, Supply: 3},
	"Snail":        {Mythium: 5,  Supply: 1},
	"Lizard":       {Mythium: 10, Supply: 2},
	"Fiend":        {Mythium: 15, Supply: 1},
	"Brute":        {Mythium: 30, Supply: 3},
	"Dragon Turtle": {Mythium: 60, Supply: 3},
	"Hermit":       {Mythium: 50, Supply: 2},
	"Dino":         {Mythium: 20, Supply: 2},
	"Safety Mole":  {Mythium: 5,  Supply: 1},
	"Drake":        {Mythium: 25, Supply: 2},
	"Mimic":        {Mythium: 40, Supply: 3},
	"Pack Leader":  {Mythium: 15, Supply: 2},
	"Centaur":      {Mythium: 10, Supply: 1},
	"Four Eyes":    {Mythium: 20, Supply: 2},
	"Shaman":       {Mythium: 15, Supply: 1},
	"Ghost Knight": {Mythium: 35, Supply: 3},
	"Kraken":       {Mythium: 55, Supply: 4},
	"Ogre":         {Mythium: 60, Supply: 4},
	"Imp":          {Mythium: 10, Supply: 1},
	"Needler":      {Mythium: 20, Supply: 2},
	"Cannoneer":    {Mythium: 15, Supply: 2},
	"Robo":         {Mythium: 40, Supply: 4},
	"Honey Bear":   {Mythium: 20, Supply: 2},
}

var FighterAttack = map[string]AttackType{
	"Proton":             AtkPierce,
	"Seadragon":           AtkMagic,
	"Angler":             AtkPierce,
	"Warp Wing":          AtkPierce,
	"Fire Lord":          AtkMagic,
	"Harpy":              AtkPierce,
	"Treant":             AtkNormal,
	"Gatling Gun":        AtkPierce,
	"Eggsack":            AtkNormal,
	"Butcher":            AtkNormal,
	"Head Chef":          AtkNormal,
	"Nightmare":          AtkPierce,
	"Doppelganger":       AtkPierce,
	"Lord Of Death":      AtkMagic,
	"Hades":              AtkMagic,
	"Old Hound":          AtkNormal,
	"Chaos Hound":        AtkNormal,
	"Diabolic":           AtkPierce,
	"Bone Warrior":       AtkNormal,
	"Bone Crusher":       AtkNormal,
	"Fire Archer":        AtkPierce,
	"Dark Mage":          AtkMagic,
	"Gargoyle":           AtkMagic,
	"Green Devil":        AtkMagic,
	"Gateguard":          AtkNormal,
	"Harbinger":          AtkNormal,
	"Bazooka":            AtkNormal,
	"Pyro":               AtkMagic,
	"Zeus":               AtkNormal,
	"APS":                AtkPierce,
	"MPS":                AtkPierce,
	"Tempest":            AtkPierce,
	"Leviathan":          AtkPierce,
	"Berserker":          AtkPierce,
	"Fatalizer":          AtkPierce,
	"Millennium":         AtkNormal,
	"Doomsday Machine":    AtkNormal,
	"Reactor":            AtkNormal,
	"Buzz":               AtkPierce,
	"Consort":            AtkPierce,
	"Ranger":             AtkPierce,
	"Daphne":             AtkPierce,
	"Wileshroom":         AtkNormal,
	"Canopie":            AtkNormal,
	"Honeyflower":        AtkMagic,
	"Deathcap":           AtkMagic,
	"Antler":             AtkNormal,
	"Whitemane":          AtkNormal,
	"Banana Bunk":        AtkPierce,
	"Banana Haven":       AtkPierce,
	"Disciple":           AtkPierce,
	"Starcaller":         AtkPierce,
	"Fire Elemental":     AtkMagic,
	"Fenix":              AtkMagic,
	"Peewee":             AtkPierce,
	"Veteran":            AtkPierce,
	"Aqua Spirit":        AtkPierce,
	"Rogue Wave":         AtkPierce,
	"Windhawk":           AtkNormal,
	"Violet":             AtkNormal,
	"Mudman":             AtkNormal,
	"Golem":              AtkNormal,
	"Pollywog":           AtkMagic,
	"Yozora":             AtkMagic,
	"Masked Spirit":      AtkPierce,
	"Sacred Steed":       AtkNormal,
	"Elite Archer":       AtkPierce,
	"Mr Brewpot":         AtkMagic,
	"Looter":             AtkPierce,
	"Howler":             AtkPierce,
	"Grarl":              AtkNormal,
	"Pulsebot":           AtkMagic,
	"Great Boar":         AtkNormal,
	"Golden Buckler":     AtkNormal,
	"Priestess Of The Abyss": AtkMagic,
	"White Mage":         AtkMagic,
	"Skybot":             AtkPierce,
	"Giga Annihilator":   AtkNormal,
	"Undead Dragon":      AtkMagic,
	"Sand Badger":        AtkNormal,
	"Pegasus":            AtkPierce,
	"Sunfang":            AtkPierce,
	"Deepcoiler":         AtkMagic,
	"Nightcrawler":       AtkPierce,
	"Tethered Soul":      AtkMagic,
	"Desert Pilgrim":     AtkNormal,
	"Soul Gate":          AtkMagic,
	"Eternal Wanderer":   AtkMagic,
	"Holy Avenger":       AtkNormal,
}

var MercAttack = map[string]AttackType{
	"Witch":       AtkMagic,
	"Siege Ram":  AtkNormal,
	"Snail":       AtkNormal,
	"Lizard":      AtkPierce,
	"Fiend":       AtkPierce,
	"Brute":       AtkNormal,
	"Dragon Turtle": AtkMagic,
	"Hermit":      AtkMagic,
	"Dino":        AtkPierce,
	"Safety Mole": AtkNormal,
	"Drake":       AtkMagic,
	"Mimic":       AtkPierce,
	"Pack Leader": AtkNormal,
	"Centaur":     AtkPierce,
	"Four Eyes":   AtkMagic,
	"Shaman":      AtkMagic,
	"Ghost Knight": AtkNormal,
	"Kraken":      AtkNormal,
	"Ogre":        AtkNormal,
	"Imp":         AtkPierce,
	"Needler":     AtkPierce,
	"Cannoneer":   AtkNormal,
	"Robo":        AtkPierce,
	"Honey Bear":  AtkMagic,
}

func GetFighterAttack(name string) (AttackType, bool) {
	at, ok := FighterAttack[name]
	return at, ok
}

func GetMercAttack(name string) (AttackType, bool) {
	at, ok := MercAttack[name]
	return at, ok
}

func GetFighterCost(name string) (UnitCost, bool) {
	c, ok := FighterCosts[name]
	return c, ok
}

func GetMercCost(name string) (UnitCost, bool) {
	c, ok := MercCosts[name]
	return c, ok
}

type ArmorType int

const (
	ArmLight    ArmorType = iota
	ArmMedium
	ArmHeavy
	ArmFortified
)

// DamageMultiplier returns the effectiveness multiplier of given attack vs armor.
func DamageMultiplier(atk AttackType, armor ArmorType) float64 {
	chart := map[AttackType]map[ArmorType]float64{
		AtkNormal: {ArmLight: 0.80, ArmMedium: 0.90, ArmHeavy: 1.15, ArmFortified: 1.15},
		AtkPierce: {ArmLight: 1.20, ArmMedium: 0.85, ArmHeavy: 1.15, ArmFortified: 0.80},
		AtkMagic:  {ArmLight: 1.00, ArmMedium: 1.25, ArmHeavy: 0.75, ArmFortified: 1.05},
	}
	return chart[atk][armor]
}

func ParseArmor(s string) ArmorType {
	switch s {
	case "light":
		return ArmLight
	case "medium":
		return ArmMedium
	case "heavy":
		return ArmHeavy
	case "fortified":
		return ArmFortified
	default:
		return ArmMedium
	}
}

// BestAttackAgainst returns which attack type is most effective vs the given armor.
func BestAttackAgainst(armor ArmorType) AttackType {
	best := AtkNormal
	bestMul := 0.0
	for _, atk := range []AttackType{AtkNormal, AtkPierce, AtkMagic} {
		m := DamageMultiplier(atk, armor)
		if m > bestMul {
			bestMul = m
			best = atk
		}
	}
	return best
}

func (a AttackType) String() string {
	switch a {
	case AtkNormal:
		return "Normal"
	case AtkPierce:
		return "Pierce"
	case AtkMagic:
		return "Magic"
	}
	return ""
}
