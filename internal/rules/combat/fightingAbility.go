package rules

func CalculateFightingAbility(class string, level int64) int64 {
	switch class {
	case "Fighter", "Barbarian", "Berserker", "Ranger", "Warlock", "Huntsman", "Paladin", "Cataphract":
		return level
	case "Magician", "Cryomancer", "Illusionist", "Necromancer", "Pyromancer", "Witch", "Priest":
		return level / 2
	case "Cleric", "Bard", "Thief", "Legerdemainist", "Druid", "Purloiner", "Assassin", "Scout":
		if level < 1 {
			return 0
		}
		if level > 12 {
			return 8
		}
		return (level + 3) / 2
	case "Monk":
		return level - 1
	case "Shaman":
		if level < 3 {
			return 0
		}
		if level > 12 {
			return 7
		}
		return (level + 1) / 2
	}
	return 0
}

func GetTargetNumber(fightingAbility int64, armorClass int64) int64 {
	// Validate inputs
	if fightingAbility < 0 || fightingAbility > 12 {
		return 0
	}
	if armorClass < -9 || armorClass > 9 {
		return 0
	}

	// The base target number for AC 9 based on FA
	baseTarget := 11 - fightingAbility

	// Each step of AC below 9 adds 1 to the target number
	acAdjustment := 9 - armorClass

	return baseTarget + acAdjustment
}
