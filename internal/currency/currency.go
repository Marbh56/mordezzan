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

// AddToPurse adds an amount of a specific denomination to the purse,
// automatically converting to the most efficient mix of coins
func AddToPurse(p *Purse, amount int64, denom Denomination) {
	// Convert the amount to copper pieces (the base unit)
	copperAmount, _ := Convert(amount, denom, CopperPieces)

	// Start with the highest denomination and work down
	if copperAmount >= 500 {
		p.PlatinumPieces += copperAmount / 500
		copperAmount = copperAmount % 500
	}

	if copperAmount >= 100 {
		p.GoldPieces += copperAmount / 100
		copperAmount = copperAmount % 100
	}

	if copperAmount >= 50 {
		p.ElectrumPieces += copperAmount / 50
		copperAmount = copperAmount % 50
	}

	if copperAmount >= 10 {
		p.SilverPieces += copperAmount / 10
		copperAmount = copperAmount % 10
	}

	p.CopperPieces += copperAmount
}

// RemoveFromPurse removes an amount of a specific denomination from the purse
// Returns true if successful, false if there aren't enough coins
func RemoveFromPurse(p *Purse, amount int64, denom Denomination) bool {
	// Calculate the total copper value in the purse
	totalCopper := (p.PlatinumPieces * 500) +
		(p.GoldPieces * 100) +
		(p.ElectrumPieces * 50) +
		(p.SilverPieces * 10) +
		p.CopperPieces

	// Calculate copper value of the amount to remove
	removalCopper, _ := Convert(amount, denom, CopperPieces)

	// Check if there's enough currency
	if removalCopper > totalCopper {
		return false
	}

	// Calculate what remains after removal
	remaining := totalCopper - removalCopper

	// Reset the purse
	p.PlatinumPieces = 0
	p.GoldPieces = 0
	p.ElectrumPieces = 0
	p.SilverPieces = 0
	p.CopperPieces = 0

	// Refill the purse with the remaining amount in the most efficient way
	AddToPurse(p, remaining, CopperPieces)

	return true
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
