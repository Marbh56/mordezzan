package combat

type MasteryLevel string

const (
	MasteryNone     MasteryLevel = ""
	MasteryMastered MasteryLevel = "mastered"
	MasteryGrand    MasteryLevel = "grand_mastery"
)

// AttackRate represents attacks per round progression
type AttackRate string

const (
	Rate_1_2 AttackRate = "1/2" // one attack every other round
	Rate_1_1 AttackRate = "1/1" // one attack per round
	Rate_3_2 AttackRate = "3/2" // three attacks over two rounds
	Rate_2_1 AttackRate = "2/1" // two attacks per round
	Rate_5_2 AttackRate = "5/2" // five attacks over two rounds
	Rate_3_1 AttackRate = "3/1" // three attacks per round
)

// WeaponMasteryModifiers contains all bonuses granted by weapon mastery
type WeaponMasteryModifiers struct {
	ToHitBonus  int        // Bonus to hit
	DamageBonus int        // Bonus to damage
	AttackRate  AttackRate // Improved attack rate
}

// progressionOrder defines the sequence of attack rate improvements
var progressionOrder = []AttackRate{
	Rate_1_2,
	Rate_1_1,
	Rate_3_2,
	Rate_2_1,
	Rate_5_2,
	Rate_3_1,
}

// GetWeaponMasteryModifiers returns all modifiers for a given mastery level
func GetWeaponMasteryModifiers(baseAttackRate AttackRate, masteryLevel MasteryLevel) WeaponMasteryModifiers {
	mods := WeaponMasteryModifiers{
		AttackRate: baseAttackRate,
	}

	switch masteryLevel {
	case MasteryMastered:
		mods.ToHitBonus = 1
		mods.DamageBonus = 1
		mods.AttackRate = getImprovedAttackRate(baseAttackRate, masteryLevel)
	case MasteryGrand:
		mods.ToHitBonus = 2
		mods.DamageBonus = 2
		mods.AttackRate = getImprovedAttackRate(baseAttackRate, masteryLevel)
	}

	return mods
}

// getImprovedAttackRate returns the improved attack rate based on mastery level
func getImprovedAttackRate(baseRate AttackRate, masteryLevel MasteryLevel) AttackRate {
	if masteryLevel == MasteryNone {
		return baseRate
	}

	// Find current position in progression
	currentIndex := -1
	for i, rate := range progressionOrder {
		if rate == baseRate {
			currentIndex = i
			break
		}
	}

	// If base rate not found in progression or already at max, return unchanged
	if currentIndex == -1 || currentIndex == len(progressionOrder)-1 {
		return baseRate
	}

	// For regular mastery, advance one step
	if masteryLevel == MasteryMastered {
		return progressionOrder[currentIndex+1]
	}

	// For grand mastery, advance two steps if possible
	if masteryLevel == MasteryGrand && currentIndex < len(progressionOrder)-2 {
		return progressionOrder[currentIndex+2]
	}

	// If can't advance two steps for grand mastery, advance one step
	if masteryLevel == MasteryGrand {
		return progressionOrder[currentIndex+1]
	}

	return baseRate
}

// GetAvailableMasterySlots returns how many weapon masteries a fighter can have
func GetAvailableMasterySlots(level int64) int {
	if level < 1 {
		return 0
	}

	// Base 2 slots at level 1
	slots := 2

	// Additional slots at levels 4, 8, and 12
	if level >= 4 {
		slots++
	}
	if level >= 8 {
		slots++
	}
	if level >= 12 {
		slots++
	}

	return slots
}

// ParseAttackRate converts a string to an AttackRate, returning a default if invalid
func ParseAttackRate(s string) AttackRate {
	switch s {
	case "1/2":
		return Rate_1_2
	case "1/1":
		return Rate_1_1
	case "3/2":
		return Rate_3_2
	case "2/1":
		return Rate_2_1
	case "5/2":
		return Rate_5_2
	case "3/1":
		return Rate_3_1
	default:
		return Rate_1_1 // Default to 1/1 if unknown
	}
}

// String returns the string representation of an AttackRate
func (ar AttackRate) String() string {
	return string(ar)
}
