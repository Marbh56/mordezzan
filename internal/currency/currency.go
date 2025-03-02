package currency

import (
	"fmt"
	"strings"
)

// CoinsPerPound represents how many coins make up one pound of weight
const CoinsPerPound = 50.0

// Denomination represents a currency type
type Denomination string

const (
	PlatinumPieces Denomination = "pp"
	GoldPieces     Denomination = "gp"
	ElectrumPieces Denomination = "ep"
	SilverPieces   Denomination = "sp"
	CopperPieces   Denomination = "cp"
)

// Purse represents a character's currency holdings
type Purse struct {
	PlatinumPieces int64
	GoldPieces     int64
	ElectrumPieces int64
	SilverPieces   int64
	CopperPieces   int64
}

// Currency conversion rates
var conversionTable = map[Denomination]map[Denomination]int64{
	PlatinumPieces: {
		PlatinumPieces: 1,
		GoldPieces:     5,
		ElectrumPieces: 10,
		SilverPieces:   50,
		CopperPieces:   500,
	},
	GoldPieces: {
		PlatinumPieces: 0,
		GoldPieces:     1,
		ElectrumPieces: 2,
		SilverPieces:   10,
		CopperPieces:   100,
	},
	ElectrumPieces: {
		PlatinumPieces: 0,
		GoldPieces:     0,
		ElectrumPieces: 1,
		SilverPieces:   5,
		CopperPieces:   50,
	},
	SilverPieces: {
		PlatinumPieces: 0,
		GoldPieces:     0,
		ElectrumPieces: 0,
		SilverPieces:   1,
		CopperPieces:   10,
	},
	CopperPieces: {
		PlatinumPieces: 0,
		GoldPieces:     0,
		ElectrumPieces: 0,
		SilverPieces:   0,
		CopperPieces:   1,
	},
}

// Convert converts an amount of one denomination to another
func Convert(amount int64, from, to Denomination) (converted, remainder int64) {
	if from == to {
		return amount, 0
	}

	rate := conversionTable[from][to]
	if rate == 0 {
		rate = conversionTable[to][from]
		if rate == 0 {
			return 0, amount
		}
		converted = amount / rate
		remainder = amount % rate
		return converted, remainder
	}

	return amount * rate, 0
}

// GetTotalCoins gets the total number of coins in the purse
func GetTotalCoins(p *Purse) int64 {
	return p.PlatinumPieces + p.GoldPieces + p.ElectrumPieces + p.SilverPieces + p.CopperPieces
}

// GetTotalWeight calculates the total weight of all coins in pounds
func GetTotalWeight(p *Purse) float64 {
	totalCoins := GetTotalCoins(p)
	return float64(totalCoins) / CoinsPerPound
}

func AddToPurse(p *Purse, amount int64, denom Denomination) {
	switch denom {
	case PlatinumPieces:
		p.PlatinumPieces += amount
	case GoldPieces:
		p.GoldPieces += amount
	case ElectrumPieces:
		p.ElectrumPieces += amount
	case SilverPieces:
		p.SilverPieces += amount
	case CopperPieces:
		p.CopperPieces += amount
	}
}

func RemoveFromPurse(p *Purse, amount int64, denom Denomination) bool {
	switch denom {
	case PlatinumPieces:
		if p.PlatinumPieces >= amount {
			p.PlatinumPieces -= amount
			return true
		}
		// Not enough platinum
		return false

	case GoldPieces:
		if p.GoldPieces >= amount {
			p.GoldPieces -= amount
			return true
		}
		// Not enough gold, try converting from platinum
		ppNeeded := (amount - p.GoldPieces + 4) / 5 // Ceiling division
		if p.PlatinumPieces >= ppNeeded {
			p.PlatinumPieces -= ppNeeded
			p.GoldPieces += ppNeeded*5 - amount
			return true
		}
		return false

	case ElectrumPieces:
		if p.ElectrumPieces >= amount {
			p.ElectrumPieces -= amount
			return true
		}
		// Try converting from higher denominations
		// First check if we have gold pieces
		if p.GoldPieces > 0 {
			gpNeeded := (amount - p.ElectrumPieces + 1) / 2 // Ceiling division
			if p.GoldPieces >= gpNeeded {
				p.GoldPieces -= gpNeeded
				p.ElectrumPieces += gpNeeded*2 - amount
				return true
			}
		}
		// Not enough direct conversion possible
		return false

	case SilverPieces:
		if p.SilverPieces >= amount {
			p.SilverPieces -= amount
			return true
		}
		// Try from electrum first
		if p.ElectrumPieces > 0 {
			epNeeded := (amount - p.SilverPieces + 4) / 5 // Ceiling division
			if p.ElectrumPieces >= epNeeded {
				p.ElectrumPieces -= epNeeded
				p.SilverPieces += epNeeded*5 - amount
				return true
			}
		}
		// Not enough direct conversion possible
		return false

	case CopperPieces:
		if p.CopperPieces >= amount {
			p.CopperPieces -= amount
			return true
		}
		// Try from silver first
		if p.SilverPieces > 0 {
			spNeeded := (amount - p.CopperPieces + 9) / 10 // Ceiling division
			if p.SilverPieces >= spNeeded {
				p.SilverPieces -= spNeeded
				p.CopperPieces += spNeeded*10 - amount
				return true
			}
		}
		// Not enough direct conversion possible
		return false
	}

	return false
}

// FormatPurse formats the purse contents as a readable string
func FormatPurse(p *Purse) string {
	var parts []string

	if p.PlatinumPieces > 0 {
		parts = append(parts, fmt.Sprintf("%d pp", p.PlatinumPieces))
	}
	if p.GoldPieces > 0 {
		parts = append(parts, fmt.Sprintf("%d gp", p.GoldPieces))
	}
	if p.ElectrumPieces > 0 {
		parts = append(parts, fmt.Sprintf("%d ep", p.ElectrumPieces))
	}
	if p.SilverPieces > 0 {
		parts = append(parts, fmt.Sprintf("%d sp", p.SilverPieces))
	}
	if p.CopperPieces > 0 {
		parts = append(parts, fmt.Sprintf("%d cp", p.CopperPieces))
	}

	if len(parts) == 0 {
		return "0 cp"
	}

	return strings.Join(parts, ", ")
}
