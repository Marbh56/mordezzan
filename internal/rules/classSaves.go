package rules

// Calculate Bonus Experience Points

// SavingThrowModifiers contains the modifiers for each type of saving throw
type SavingThrowModifiers struct {
	Death          int64
	Transformation int64
	Device         int64
	Avoidance      int64
	Sorcery        int64
}

// GetSavingThrowModifiers returns the saving throw modifiers for a fighter
func GetSavingThrowModifiers(class string) SavingThrowModifiers {
	var mods SavingThrowModifiers

	switch class {
	case "Fighter", "Huntsman", "Ranger":
		mods.Death = -2
		mods.Transformation = -2

	case "Magician", "Cryomancer", "Illusionist", "Pyromancer":
		mods.Device = -2
		mods.Sorcery = -2

	case "Cleric", "Necromancer", "Druid", "Priest", "Shaman":
		mods.Death = -2
		mods.Sorcery = -2

	case "Theif", "Assassin", "Bard", "Scout":
		mods.Device = -2
		mods.Avoidance = -2

	case "Barbarian", "Berserker", "Paladin":
		mods.Death = -2
		mods.Device = -2
		mods.Sorcery = -2
		mods.Avoidance = -2
		mods.Transformation = -2

	case "Cataphract":
		mods.Death = -2
		mods.Transformation = -2

	case "Warlock", "Witch", "Runegraver":
		mods.Transformation = -2
		mods.Sorcery = -2

	case "Monk":
		mods.Transformation = -2
		mods.Avoidance = -2

	case "Legerdemainist", "Purloiner":
		mods.Avoidance = -2
		mods.Sorcery = -2
	}
	return mods
}
