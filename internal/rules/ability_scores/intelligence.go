package ability_scores

// IntelligenceModifiers contains all the modifiers and bonuses for a given Intelligence score
type IntelligenceModifiers struct {
	Score           int64 `json:"score"`
	Languages       int   `json:"languages"`         // Additional languages known
	BonusSpellLevel int   `json:"bonus_spell_level"` // Highest level of bonus spell granted
	ChanceToLearn   int   `json:"chance_to_learn"`   // Percentage chance to learn new spells
	IsLiterate      bool  `json:"is_literate"`       // Whether the character can read and write
}

// CalculateIntelligenceModifiers returns all intelligence-based modifiers for a given score
func CalculateIntelligenceModifiers(intelligence int64) IntelligenceModifiers {
	mods := IntelligenceModifiers{Score: intelligence}

	switch {
	case intelligence == 3:
		mods.Languages = 0
		mods.BonusSpellLevel = 0
		mods.ChanceToLearn = 0
		mods.IsLiterate = false

	case intelligence >= 4 && intelligence <= 6:
		mods.Languages = 0
		mods.BonusSpellLevel = 0
		mods.ChanceToLearn = 0
		mods.IsLiterate = false

	case intelligence >= 7 && intelligence <= 8:
		mods.Languages = 0
		mods.BonusSpellLevel = 0
		mods.ChanceToLearn = 0
		mods.IsLiterate = true

	case intelligence >= 9 && intelligence <= 12:
		mods.Languages = 0
		mods.BonusSpellLevel = 0
		mods.ChanceToLearn = 50
		mods.IsLiterate = true

	case intelligence >= 13 && intelligence <= 14:
		mods.Languages = 1
		mods.BonusSpellLevel = 1 // One level 1 spell
		mods.ChanceToLearn = 65
		mods.IsLiterate = true

	case intelligence >= 15 && intelligence <= 16:
		mods.Languages = 1
		mods.BonusSpellLevel = 2 // One level 2 spell
		mods.ChanceToLearn = 75
		mods.IsLiterate = true

	case intelligence == 17:
		mods.Languages = 2
		mods.BonusSpellLevel = 3 // One level 3 spell
		mods.ChanceToLearn = 85
		mods.IsLiterate = true

	case intelligence == 18:
		mods.Languages = 3
		mods.BonusSpellLevel = 4 // One level 4 spell
		mods.ChanceToLearn = 95
		mods.IsLiterate = true
	}

	return mods
}
