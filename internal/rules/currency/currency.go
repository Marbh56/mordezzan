package currency

import (
	"fmt"
	"strings"
)

const CoinsPerPound = 100.0

type Denomination string

const (
    PlatinumPieces Denomination = "pp"
    GoldPieces     Denomination = "gp"
    ElectrumPieces Denomination = "ep"
    SilverPieces   Denomination = "sp"
    CopperPieces   Denomination = "cp"
)

type Purse struct {
    PlatinumPieces int64
    GoldPieces     int64
    ElectrumPieces int64
    SilverPieces   int64
    CopperPieces   int64
}

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

func GetTotalCoins(p *Purse) int64 {
    return p.PlatinumPieces + p.GoldPieces + p.ElectrumPieces + p.SilverPieces + p.CopperPieces
}

func GetTotalWeight(p *Purse) float64 {
    totalCoins := GetTotalCoins(p)
    return float64(totalCoins) / CoinsPerPound
}

func AddToPurse(p *Purse, amount int64, denom Denomination) {
    copperAmount, _ := Convert(amount, denom, CopperPieces)

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

func RemoveFromPurse(p *Purse, amount int64, denom Denomination) bool {
    totalCopper := (p.PlatinumPieces * 500) +
                  (p.GoldPieces * 100) +
                  (p.ElectrumPieces * 50) +
                  (p.SilverPieces * 10) +
                  p.CopperPieces

    removalCopper, _ := Convert(amount, denom, CopperPieces)

    if removalCopper > totalCopper {
        return false
    }

    remaining := totalCopper - removalCopper

    p.PlatinumPieces = 0
    p.GoldPieces = 0
    p.ElectrumPieces = 0
    p.SilverPieces = 0
    p.CopperPieces = 0

    AddToPurse(p, remaining, CopperPieces)

    return true
}

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
