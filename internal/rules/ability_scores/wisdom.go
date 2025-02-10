package rules

// WisdomModifiers contains all modifiers and chances for a given Wisdom score
type WisdomModifiers struct {
	Score               int64 `json:"score"`
	WillpowerAdj        int   `json:"willpower_adj"`         // Willpower adjustment
	BonusSpellLevel     int   `json:"bonus_spell_level"`     // Highest level of bonus spell granted
	SpellLearningChance int   `json:"spell_learning_chance"` // Percentage chance to learn new spells
}

// CalculateWisdomModifiers returns all wisdom-based modifiers for a given score
func CalculateWisdomModifiers(wisdom int64) WisdomModifiers {
	mods := WisdomModifiers{Score: wisdom}

	switch {
	case wisdom == 3:
		mods.WillpowerAdj = -2
		mods.BonusSpellLevel = 0
		mods.SpellLearningChance = 0

	case wisdom >= 4 && wisdom <= 6:
		mods.WillpowerAdj = -1
		mods.BonusSpellLevel = 0
		mods.SpellLearningChance = 0

	case wisdom >= 7 && wisdom <= 8:
		mods.WillpowerAdj = 0
		mods.BonusSpellLevel = 0
		mods.SpellLearningChance = 0

	case wisdom >= 9 && wisdom <= 12:
		mods.WillpowerAdj = 0
		mods.BonusSpellLevel = 0
		mods.SpellLearningChance = 50

	case wisdom >= 13 && wisdom <= 14:
		mods.WillpowerAdj = 0
		mods.BonusSpellLevel = 1 // One level 1 spell
		mods.SpellLearningChance = 65

	case wisdom >= 15 && wisdom <= 16:
		mods.WillpowerAdj = 1
		mods.BonusSpellLevel = 2 // One level 2 spell
		mods.SpellLearningChance = 75

	case wisdom == 17:
		mods.WillpowerAdj = 1
		mods.BonusSpellLevel = 3 // One level 3 spell
		mods.SpellLearningChance = 85

	case wisdom == 18:
		mods.WillpowerAdj = 2
		mods.BonusSpellLevel = 4 // One level 4 spell
		mods.SpellLearningChance = 95
	}

	return mods
}

// GetBonusSpells returns a map of spell levels and how many bonus spells are granted
func (w WisdomModifiers) GetBonusSpells() map[int]int {
	bonusSpells := make(map[int]int)

	// No bonus spells if wisdom doesn't grant any
	if w.BonusSpellLevel == 0 {
		return bonusSpells
	}

	// Add one bonus spell for each level up to the maximum bonus level
	for level := 1; level <= w.BonusSpellLevel; level++ {
		bonusSpells[level] = 1
	}

	return bonusSpells
}
