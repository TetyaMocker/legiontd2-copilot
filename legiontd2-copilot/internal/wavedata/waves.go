package wavedata

type EnemyArmor string

const (
	ArmorLight    EnemyArmor = "light"
	ArmorMedium   EnemyArmor = "medium"
	ArmorHeavy    EnemyArmor = "heavy"
	ArmorFortified EnemyArmor = "fortified"
)

type EnemyAttack string

const (
	AttackNormal  EnemyAttack = "normal"
	AttackPierce  EnemyAttack = "pierce"
	AttackMagic   EnemyAttack = "magic"
	AttackChaos   EnemyAttack = "chaos"
)

type Wave struct {
	Number       int
	Name        string
	EnemyName   string
	ArmorType   EnemyArmor
	AttackType  EnemyAttack
	Amount      int
	Amount2     int
	BossName    string
	BossArmor   EnemyArmor
	BossAttack  EnemyAttack
	PrepareTime int
	RecValue    int
	TotalReward int
}

var WavesByNumber = map[int]Wave{
	1:  {1, "Crabs", "Crab", ArmorFortified, AttackPierce, 12, 0, "", "", "", 90, 150, 72},
	2:  {2, "Wales", "Wale", ArmorHeavy, AttackNormal, 12, 0, "", "", "", 25, 150, 84},
	3:  {3, "Hoppers", "Hopper", ArmorMedium, AttackMagic, 18, 0, "", "", "", 26, 215, 90},
	4:  {4, "Flying Chickens", "Flying Chicken", ArmorLight, AttackNormal, 12, 0, "", "", "", 27, 270, 96},
	5:  {5, "Scorpions", "Scorpion", ArmorMedium, AttackPierce, 8, 1, "Scorpion King", ArmorMedium, AttackPierce, 28, 335, 108},
	6:  {6, "Rockos", "Rocko", ArmorFortified, AttackNormal, 6, 0, "", "", "", 29, 445, 114},
	7:  {7, "Sludges", "Sludge", ArmorHeavy, AttackMagic, 10, 0, "", "", "", 30, 570, 120},
	8:  {8, "Kobras", "Kobra", ArmorLight, AttackMagic, 12, 0, "", "", "", 31, 700, 132},
	9:  {9, "Carapaces", "Carapace", ArmorFortified, AttackPierce, 12, 0, "", "", "", 32, 820, 144},
	10: {10, "Granddaddy", "Granddaddy", ArmorHeavy, AttackNormal, 1, 0, "", "", "", 33, 1075, 150},
	11: {11, "Quill Shooters", "Quill Shooter", ArmorMedium, AttackPierce, 12, 0, "", "", "", 45, 1380, 156},
	12: {12, "Mantises", "Mantis", ArmorLight, AttackPierce, 12, 0, "", "", "", 35, 1570, 168},
	13: {13, "Drill Golems", "Drill Golem", ArmorFortified, AttackNormal, 6, 0, "", "", "", 36, 1920, 180},
	14: {14, "Killer Slugs", "Killer Slug", ArmorHeavy, AttackMagic, 12, 0, "", "", "", 37, 2200, 192},
	15: {15, "Quadrapuses", "Quadrapus", ArmorMedium, AttackMagic, 8, 1, "Giant Quadrapus", ArmorMedium, AttackMagic, 38, 2800, 204},
	16: {16, "Cardinals", "Cardinal", ArmorLight, AttackNormal, 18, 0, "", "", "", 39, 3350, 216},
	17: {17, "Metal Dragons", "Metal Dragon", ArmorHeavy, AttackPierce, 12, 0, "", "", "", 40, 4050, 228},
	18: {18, "Wale Chiefs", "Wale Chief", ArmorMedium, AttackNormal, 6, 0, "", "", "", 41, 5000, 252},
	19: {19, "Dire Toads", "Dire Toad", ArmorLight, AttackPierce, 12, 0, "", "", "", 42, 6100, 276},
	20: {20, "Maccabeus", "Maccabeus", ArmorFortified, AttackMagic, 1, 0, "", "", "", 43, 7350, 300},
	21: {21, "Legion Lord", "Legion Lord", ArmorLight, AttackChaos, 8, 1, "Legion King", ArmorLight, AttackChaos, 44, 9000, 360},
}

func GetWave(num int) (Wave, bool) {
	w, ok := WavesByNumber[num]
	return w, ok
}
