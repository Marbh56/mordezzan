package rules

type EncumbranceThresholds struct {
	Score               int64 `json:"score"`
	BaseEncumbered      int   `json:"base_encumbered"`       // -10 MV, -1 AC
	BaseHeavyEncumbered int   `json:"base_heavy_encumbered"` // -20 MV, -2 AC
	MaximumCapacity     int   `json:"maximum_capacity"`      // Cannot carry more than this
}

func CalculateEncumbranceThresholds(strength, constitution int64) EncumbranceThresholds {
	baseThresholds := EncumbranceThresholds{
		BaseEncumbered:      75,
		BaseHeavyEncumbered: 150,
		MaximumCapacity:     300, // Base maximum capacity
	}

	// Calculate strength modifier (in pounds)
	strMod := 0
	maxMod := 0
	switch {
	case strength <= 6:
		strMod = -25
		maxMod = -100
	case strength >= 7 && strength <= 8:
		strMod = -15
		maxMod = -50
	case strength >= 13 && strength <= 14:
		strMod = 15
		maxMod = 50
	case strength >= 15 && strength <= 16:
		strMod = 25
		maxMod = 100
	case strength == 17:
		strMod = 35
		maxMod = 150
	case strength == 18:
		strMod = 50
		maxMod = 200
	}

	// Calculate constitution modifier (in pounds)
	conMod := 0
	conMaxMod := 0
	switch {
	case constitution <= 6:
		conMod = -10
		conMaxMod = -25
	case constitution >= 7 && constitution <= 8:
		conMod = -5
		conMaxMod = -15
	case constitution >= 13 && constitution <= 14:
		conMod = 5
		conMaxMod = 15
	case constitution >= 15 && constitution <= 16:
		conMod = 10
		conMaxMod = 25
	case constitution >= 17:
		conMod = 15
		conMaxMod = 35
	}

	// Apply modifiers
	baseThresholds.BaseEncumbered += strMod + conMod
	baseThresholds.BaseHeavyEncumbered += (strMod + conMod) * 2
	baseThresholds.MaximumCapacity += maxMod + conMaxMod

	// Ensure minimum thresholds
	if baseThresholds.BaseEncumbered < 40 {
		baseThresholds.BaseEncumbered = 40
	}
	if baseThresholds.BaseHeavyEncumbered < 60 {
		baseThresholds.BaseHeavyEncumbered = 60
	}
	if baseThresholds.MaximumCapacity < 100 {
		baseThresholds.MaximumCapacity = 100
	}

	return baseThresholds
}
