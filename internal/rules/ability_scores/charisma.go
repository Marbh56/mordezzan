package ability_scores

// CharismaModifiers contains all modifiers and limitations for a given Charisma score
type CharismaModifiers struct {
	Score              int64 `json:"score"`
	ReactionLoyaltyAdj int   `json:"reaction_loyalty_adj"` // Reaction and loyalty adjustment
	MaxHenchmen        int   `json:"max_henchmen"`         // Maximum number of henchmen
	UndeadTurningAdj   int   `json:"undead_turning_adj"`   // Undead turning adjustment
}

// CalculateCharismaModifiers returns all charisma-based modifiers for a given score
func CalculateCharismaModifiers(charisma int64) CharismaModifiers {
	mods := CharismaModifiers{Score: charisma}

	switch {
	case charisma == 3:
		mods.ReactionLoyaltyAdj = -3
		mods.MaxHenchmen = 1
		mods.UndeadTurningAdj = -1

	case charisma >= 4 && charisma <= 6:
		mods.ReactionLoyaltyAdj = -2
		mods.MaxHenchmen = 2
		mods.UndeadTurningAdj = -1

	case charisma >= 7 && charisma <= 8:
		mods.ReactionLoyaltyAdj = -1
		mods.MaxHenchmen = 3
		mods.UndeadTurningAdj = 0

	case charisma >= 9 && charisma <= 12:
		mods.ReactionLoyaltyAdj = 0
		mods.MaxHenchmen = 4
		mods.UndeadTurningAdj = 0

	case charisma >= 13 && charisma <= 14:
		mods.ReactionLoyaltyAdj = 1
		mods.MaxHenchmen = 6
		mods.UndeadTurningAdj = 0

	case charisma >= 15 && charisma <= 16:
		mods.ReactionLoyaltyAdj = 1
		mods.MaxHenchmen = 8
		mods.UndeadTurningAdj = 1

	case charisma == 17:
		mods.ReactionLoyaltyAdj = 2
		mods.MaxHenchmen = 10
		mods.UndeadTurningAdj = 1

	case charisma == 18:
		mods.ReactionLoyaltyAdj = 3
		mods.MaxHenchmen = 12
		mods.UndeadTurningAdj = 1
	}

	return mods
}

// GetLoyaltyBonus returns the loyalty bonus for henchmen and retainers
func (c CharismaModifiers) GetLoyaltyBonus() int {
	return c.ReactionLoyaltyAdj
}

// GetReactionBonus returns the reaction adjustment for NPC encounters
func (c CharismaModifiers) GetReactionBonus() int {
	return c.ReactionLoyaltyAdj
}

// GetTurningBonus returns the bonus to turn undead attempts
func (c CharismaModifiers) GetTurningBonus() int {
	return c.UndeadTurningAdj
}
