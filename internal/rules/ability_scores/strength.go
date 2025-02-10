package rules

// StrengthModifiers contains all the modifiers and chances for a given Strength score
type StrengthModifiers struct {
	Score             int64  `json:"score"`
	AttackMod         int    `json:"attack_mod"`         // To-hit modifier for melee
	DamageMod         int    `json:"damage_mod"`         // Damage adjustment for melee/hurled
	TestOfStrength    string `json:"test_of_strength"`   // X:6 chance format
	ExtraordinaryFeat int    `json:"extraordinary_feat"` // Percentage chance
}

// CalculateStrengthModifiers returns all strength-based modifiers for a given score
func CalculateStrengthModifiers(strength int64) StrengthModifiers {
	mods := StrengthModifiers{Score: strength}

	switch {
	case strength == 3:
		mods.AttackMod = -2
		mods.DamageMod = -2
		mods.TestOfStrength = "1:6"
		mods.ExtraordinaryFeat = 0

	case strength >= 4 && strength <= 6:
		mods.AttackMod = -1
		mods.DamageMod = -1
		mods.TestOfStrength = "1:6"
		mods.ExtraordinaryFeat = 1

	case strength >= 7 && strength <= 8:
		mods.AttackMod = 0
		mods.DamageMod = -1
		mods.TestOfStrength = "2:6"
		mods.ExtraordinaryFeat = 2

	case strength >= 9 && strength <= 12:
		mods.AttackMod = 0
		mods.DamageMod = 0
		mods.TestOfStrength = "2:6"
		mods.ExtraordinaryFeat = 4

	case strength >= 13 && strength <= 14:
		mods.AttackMod = 0
		mods.DamageMod = 1
		mods.TestOfStrength = "3:6"
		mods.ExtraordinaryFeat = 8

	case strength >= 15 && strength <= 16:
		mods.AttackMod = 1
		mods.DamageMod = 1
		mods.TestOfStrength = "3:6"
		mods.ExtraordinaryFeat = 16

	case strength == 17:
		mods.AttackMod = 1
		mods.DamageMod = 2
		mods.TestOfStrength = "4:6"
		mods.ExtraordinaryFeat = 24

	case strength == 18:
		mods.AttackMod = 2
		mods.DamageMod = 3
		mods.TestOfStrength = "5:6"
		mods.ExtraordinaryFeat = 32
	}

	return mods
}
