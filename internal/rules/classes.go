package rules

type ClassLevel struct {
	Level       int64
	XPRequired  int64
	HitDice     string
	SavingThrow int64
	Spells      SpellSlots
}

type SpellSlots struct {
	Level1 int
	Level2 int
	Level3 int
	Level4 int
	Level5 int
	Level6 int
}

type ClassProgression struct {
	Name   string
	Levels []ClassLevel
}

var FighterProgression = ClassProgression{
	Name: "Fighter",
	Levels: []ClassLevel{
		{1, 0, "1d10", 16, SpellSlots{}},
		{2, 2000, "2d10", 16, SpellSlots{}},
		{3, 4000, "3d10", 15, SpellSlots{}},
		{4, 8000, "4d10", 15, SpellSlots{}},
		{5, 16000, "5d10", 14, SpellSlots{}},
		{6, 32000, "6d10", 14, SpellSlots{}},
		{7, 64000, "7d10", 13, SpellSlots{}},
		{8, 128000, "8d10", 13, SpellSlots{}},
		{9, 256000, "9d10", 12, SpellSlots{}},
		{10, 384000, "9d10+3", 12, SpellSlots{}},
		{11, 512000, "9d10+6", 11, SpellSlots{}},
		{12, 640000, "9d10+9", 11, SpellSlots{}},
	},
}

var MagicianProgression = ClassProgression{
	Name: "Magician",
	Levels: []ClassLevel{
		{1, 0, "1d4", 16, SpellSlots{1, 0, 0, 0, 0, 0}},
		{2, 2500, "2d4", 16, SpellSlots{2, 0, 0, 0, 0, 0}},
		{3, 5000, "3d4", 15, SpellSlots{2, 1, 0, 0, 0, 0}},
		{4, 10000, "4d4", 15, SpellSlots{3, 2, 0, 0, 0, 0}},
		{5, 20000, "5d4", 14, SpellSlots{3, 2, 1, 0, 0, 0}},
		{6, 40000, "6d4", 14, SpellSlots{4, 3, 2, 0, 0, 0}},
		{7, 80000, "7d4", 13, SpellSlots{4, 3, 2, 1, 0, 0}},
		{8, 160000, "8d4", 13, SpellSlots{4, 4, 3, 2, 0, 0}},
		{9, 320000, "9d4", 12, SpellSlots{5, 4, 3, 2, 1, 0}},
		{10, 480000, "9d4+1", 12, SpellSlots{5, 4, 4, 3, 2, 0}},
		{11, 640000, "9d4+2", 11, SpellSlots{5, 5, 4, 3, 2, 1}},
		{12, 800000, "9d4+3", 11, SpellSlots{5, 5, 4, 4, 3, 2}},
	},
}

var ClericProgression = ClassProgression{
	Name: "Cleric",
	Levels: []ClassLevel{
		{1, 0, "1d8", 16, SpellSlots{1, 0, 0, 0, 0, 0}},
		{2, 2000, "2d8", 16, SpellSlots{2, 0, 0, 0, 0, 0}},
		{3, 4000, "3d8", 15, SpellSlots{2, 1, 0, 0, 0, 0}},
		{4, 8000, "4d8", 15, SpellSlots{2, 2, 0, 0, 0, 0}},
		{5, 16000, "5d8", 14, SpellSlots{3, 2, 1, 0, 0, 0}},
		{6, 32000, "6d8", 14, SpellSlots{3, 2, 2, 0, 0, 0}},
		{7, 64000, "7d8", 13, SpellSlots{3, 3, 2, 1, 0, 0}},
		{8, 128000, "8d8", 13, SpellSlots{3, 3, 2, 2, 0, 0}},
		{9, 256000, "9d8", 12, SpellSlots{4, 3, 3, 2, 1, 0}},
		{10, 384000, "9d8+2", 12, SpellSlots{4, 3, 3, 2, 2, 0}},
		{11, 512000, "9d8+4", 11, SpellSlots{4, 4, 3, 3, 2, 1}},
		{12, 640000, "9d8+6", 11, SpellSlots{4, 4, 3, 3, 2, 2}},
	},
}

var ThiefProgression = ClassProgression{
	Name: "Thief",
	Levels: []ClassLevel{
		{1, 0, "1d6", 16, SpellSlots{}},
		{2, 1500, "2d6", 16, SpellSlots{}},
		{3, 3000, "3d6", 15, SpellSlots{}},
		{4, 6000, "4d6", 15, SpellSlots{}},
		{5, 12000, "5d6", 14, SpellSlots{}},
		{6, 24000, "6d6", 14, SpellSlots{}},
		{7, 48000, "7d6", 13, SpellSlots{}},
		{8, 96000, "8d6", 13, SpellSlots{}},
		{9, 192000, "9d6", 12, SpellSlots{}},
		{10, 288000, "9d6+2", 12, SpellSlots{}},
		{11, 384000, "9d6+4", 11, SpellSlots{}},
		{12, 480000, "9d6+6", 11, SpellSlots{}},
	},
}

func GetClassProgression(className string) ClassProgression {
	switch className {
	case "Magician":
		return MagicianProgression
	case "Fighter":
		return FighterProgression
	default:
		return ClassProgression{}
	}
}

// GetLevelForXP returns the level a character should be based on their XP
func (c ClassProgression) GetLevelForXP(xp int64) int64 {
	for i := len(c.Levels) - 1; i >= 0; i-- {
		if xp >= c.Levels[i].XPRequired {
			return c.Levels[i].Level
		}
	}
	return 1
}

// GetXPForNextLevel returns how much XP is needed for next level
func (c ClassProgression) GetXPForNextLevel(currentXP int64) int64 {
	for _, level := range c.Levels {
		if level.XPRequired > currentXP {
			return level.XPRequired - currentXP
		}
	}
	return 0 // Already at max level
}

// GetSavingThrow returns the saving throw value for a given level
func (c ClassProgression) GetSavingThrow(level int64) int64 {
	if level < 1 {
		return c.Levels[0].SavingThrow
	}
	for _, l := range c.Levels {
		if l.Level == level {
			return l.SavingThrow
		}
	}
	// If level is beyond progression, return last defined saving throw
	return c.Levels[len(c.Levels)-1].SavingThrow
}

// GetHitDice returns the hit dice string for a given level
func (c ClassProgression) GetHitDice(level int64) string {
	if level < 1 {
		return c.Levels[0].HitDice
	}
	for _, l := range c.Levels {
		if l.Level == level {
			return l.HitDice
		}
	}
	// If level is beyond progression, return last defined hit dice
	return c.Levels[len(c.Levels)-1].HitDice
}
