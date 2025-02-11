package ability_scores
// DexterityModifiers contains all modifiers and chances for a given Dexterity score
type DexterityModifiers struct {
	Score             int64  `json:"score"`
	AttackMod         int    `json:"attack_mod"`         // To-hit modifier for missile
	DefenseAdj        int    `json:"defense_adj"`        // Defense adjustment
	TestOfDexterity   string `json:"test_of_dexterity"`  // X:6 chance format
	ExtraordinaryFeat int    `json:"extraordinary_feat"` // Percentage chance
}

// CalculateDexterityModifiers returns all dexterity-based modifiers for a given score
func CalculateDexterityModifiers(dexterity int64) DexterityModifiers {
	mods := DexterityModifiers{Score: dexterity}

	switch {
	case dexterity == 3:
		mods.AttackMod = -2
		mods.DefenseAdj = -2
		mods.TestOfDexterity = "1:6"
		mods.ExtraordinaryFeat = 0

	case dexterity >= 4 && dexterity <= 6:
		mods.AttackMod = -1
		mods.DefenseAdj = -1
		mods.TestOfDexterity = "1:6"
		mods.ExtraordinaryFeat = 1

	case dexterity >= 7 && dexterity <= 8:
		mods.AttackMod = -1
		mods.DefenseAdj = 0
		mods.TestOfDexterity = "2:6"
		mods.ExtraordinaryFeat = 2

	case dexterity >= 9 && dexterity <= 12:
		mods.AttackMod = 0
		mods.DefenseAdj = 0
		mods.TestOfDexterity = "2:6"
		mods.ExtraordinaryFeat = 4

	case dexterity >= 13 && dexterity <= 14:
		mods.AttackMod = 1
		mods.DefenseAdj = 0
		mods.TestOfDexterity = "3:6"
		mods.ExtraordinaryFeat = 8

	case dexterity >= 15 && dexterity <= 16:
		mods.AttackMod = 1
		mods.DefenseAdj = 1
		mods.TestOfDexterity = "3:6"
		mods.ExtraordinaryFeat = 16

	case dexterity == 17:
		mods.AttackMod = 2
		mods.DefenseAdj = 1
		mods.TestOfDexterity = "4:6"
		mods.ExtraordinaryFeat = 24

	case dexterity == 18:
		mods.AttackMod = 3
		mods.DefenseAdj = 2
		mods.TestOfDexterity = "5:6"
		mods.ExtraordinaryFeat = 32
	}

	return mods
}
