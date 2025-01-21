package rules

// ConstitutionModifiers contains all modifiers and chances for a given Constitution score
type ConstitutionModifiers struct {
	Score             int64  `json:"score"`
	HitPointMod       int    `json:"hp_mod"`             // HP adjustment per level
	PoisonRadMod      int    `json:"poison_rad_mod"`     // Poison/Radiation adjustment
	TraumaSurvival    int    `json:"trauma_survival"`    // Percentage chance
	TestOfCon         string `json:"test_of_con"`        // X:6 chance format
	ExtraordinaryFeat int    `json:"extraordinary_feat"` // Percentage chance
}

// CalculateConstitutionModifiers returns all constitution-based modifiers for a given score
func CalculateConstitutionModifiers(constitution int64) ConstitutionModifiers {
	mods := ConstitutionModifiers{Score: constitution}

	switch {
	case constitution == 3:
		mods.HitPointMod = -1
		mods.PoisonRadMod = -2
		mods.TraumaSurvival = 45
		mods.TestOfCon = "1:6"
		mods.ExtraordinaryFeat = 0

	case constitution >= 4 && constitution <= 6:
		mods.HitPointMod = -1
		mods.PoisonRadMod = -1
		mods.TraumaSurvival = 55
		mods.TestOfCon = "1:6"
		mods.ExtraordinaryFeat = 1

	case constitution >= 7 && constitution <= 8:
		mods.HitPointMod = 0
		mods.PoisonRadMod = 0
		mods.TraumaSurvival = 65
		mods.TestOfCon = "2:6"
		mods.ExtraordinaryFeat = 2

	case constitution >= 9 && constitution <= 12:
		mods.HitPointMod = 0
		mods.PoisonRadMod = 0
		mods.TraumaSurvival = 75
		mods.TestOfCon = "2:6"
		mods.ExtraordinaryFeat = 4

	case constitution >= 13 && constitution <= 14:
		mods.HitPointMod = 1
		mods.PoisonRadMod = 0
		mods.TraumaSurvival = 80
		mods.TestOfCon = "3:6"
		mods.ExtraordinaryFeat = 8

	case constitution >= 15 && constitution <= 16:
		mods.HitPointMod = 1
		mods.PoisonRadMod = 1
		mods.TraumaSurvival = 85
		mods.TestOfCon = "3:6"
		mods.ExtraordinaryFeat = 16

	case constitution == 17:
		mods.HitPointMod = 2
		mods.PoisonRadMod = 1
		mods.TraumaSurvival = 90
		mods.TestOfCon = "4:6"
		mods.ExtraordinaryFeat = 24

	case constitution == 18:
		mods.HitPointMod = 3
		mods.PoisonRadMod = 2
		mods.TraumaSurvival = 95
		mods.TestOfCon = "5:6"
		mods.ExtraordinaryFeat = 32
	}

	return mods
}
