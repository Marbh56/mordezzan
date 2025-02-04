package currency

import (
	"fmt"
	"strings"
)

// Denomination represents a currency denomination
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

// conversionTable represents the conversion rates between currencies
// The first key is the from currency, the second key is the to currency
var conversionTable = map[Denomination]map[Denomination]int64{
    PlatinumPieces: {
        PlatinumPieces: 1,
        GoldPieces:     5,
        ElectrumPieces: 10,
        SilverPieces:   50,
        CopperPieces:   250,
    },
    GoldPieces: {
        PlatinumPieces: 0, // Requires 5 gp for 1 pp
        GoldPieces:     1,
        ElectrumPieces: 2,
        SilverPieces:   10,
        CopperPieces:   50,
    },
    ElectrumPieces: {
        PlatinumPieces: 0,
        GoldPieces:     0,
        ElectrumPieces: 1,
        SilverPieces:   5,
        CopperPieces:   25,
    },
    SilverPieces: {
        PlatinumPieces: 0,
        GoldPieces:     0,
        ElectrumPieces: 0,
        SilverPieces:   1,
        CopperPieces:   5,
    },
    CopperPieces: {
        PlatinumPieces: 0,
        GoldPieces:     0,
        ElectrumPieces: 0,
        SilverPieces:   0,
        CopperPieces:   1,
    },
}

// Convert converts an amount from one denomination to another
// Returns the converted amount and any remainder in the original denomination
func Convert(amount int64, from, to Denomination) (converted, remainder int64) {
    if from == to {
        return amount, 0
    }

    rate := conversionTable[from][to]
    if rate == 0 {
        // Handle conversion to higher denominations
        rate = conversionTable[to][from]
        if rate == 0 {
            return 0, amount // Invalid conversion
        }
        converted = amount / rate
        remainder = amount % rate
        return converted, remainder
    }

    return amount * rate, 0
}

// AddToPurse adds the specified amount of currency to a purse
// It automatically converts to the highest possible denominations
func AddToPurse(p *Purse, amount int64, denom Denomination) {
    // First, convert everything to copper pieces for easier calculation
    copperAmount, _ := Convert(amount, denom, CopperPieces)

    // Then convert to the highest possible denominations
    if copperAmount >= 250 {
        p.PlatinumPieces += copperAmount / 250
        copperAmount = copperAmount % 250
    }

    if copperAmount >= 50 {
        p.GoldPieces += copperAmount / 50
        copperAmount = copperAmount % 50
    }

    if copperAmount >= 25 {
        p.ElectrumPieces += copperAmount / 25
        copperAmount = copperAmount % 25
    }

    if copperAmount >= 5 {
        p.SilverPieces += copperAmount / 5
        copperAmount = copperAmount % 5
    }

    p.CopperPieces += copperAmount
}

// RemoveFromPurse attempts to remove the specified amount of currency from a purse
// Returns true if successful, false if insufficient funds
func RemoveFromPurse(p *Purse, amount int64, denom Denomination) bool {
    // Convert to copper pieces for comparison
    totalCopper := (p.PlatinumPieces * 250) +
                  (p.GoldPieces * 50) +
                  (p.ElectrumPieces * 25) +
                  (p.SilverPieces * 5) +
                  p.CopperPieces

    removalCopper, _ := Convert(amount, denom, CopperPieces)

    if removalCopper > totalCopper {
        return false
    }

    // Create a temporary purse to handle the remaining amount
    remaining := totalCopper - removalCopper

    // Reset the purse
    p.PlatinumPieces = 0
    p.GoldPieces = 0
    p.ElectrumPieces = 0
    p.SilverPieces = 0
    p.CopperPieces = 0

    // Add back the remaining amount
    AddToPurse(p, remaining, CopperPieces)

    return true
}

// GetTotalInCopper returns the total value of a purse in copper pieces
func GetTotalInCopper(p *Purse) int64 {
    return (p.PlatinumPieces * 250) +
           (p.GoldPieces * 50) +
           (p.ElectrumPieces * 25) +
           (p.SilverPieces * 5) +
           p.CopperPieces
}

// FormatPurse returns a formatted string representation of a purse
func FormatPurse(p *Purse) string {
    var result string
    if p.PlatinumPieces > 0 {
        result += fmt.Sprintf("%d pp ", p.PlatinumPieces)
    }
    if p.GoldPieces > 0 {
        result += fmt.Sprintf("%d gp ", p.GoldPieces)
    }
    if p.ElectrumPieces > 0 {
        result += fmt.Sprintf("%d ep ", p.ElectrumPieces)
    }
    if p.SilverPieces > 0 {
        result += fmt.Sprintf("%d sp ", p.SilverPieces)
    }
    if p.CopperPieces > 0 {
        result += fmt.Sprintf("%d cp", p.CopperPieces)
    }
    return strings.TrimSpace(result)
}
