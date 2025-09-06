package model

import (
	"sort"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type ScraperBond struct {
	Name                    string  `json:"name"`
	Investable              bool    `json:"investable"`
	AnnualInvestmentRate    string  `json:"annual_investment_rate"`
	UnitaryInvestmentValue  float64 `json:"unitary_investment_value"`
	MinimumInvestmentAmount float64 `json:"minimum_investment_amount"`
	AnnualRedemptionRate    string  `json:"annual_redemption_rate"`
	UnitaryRedemptionValue  float64 `json:"unitary_redemption_value"`
}

type ScraperCache struct {
	BondsMap  map[string]ScraperBond
	BondsList []ScraperBond
	UpdatedAt string
}

func (sc *ScraperCache) Save(investData []Invest, redeemData []Redeem) {
	bondsMap := make(map[string]ScraperBond)

	for _, red := range redeemData {
		bondsMap[strings.ToLower(red.Name)] = ScraperBond{
			Name:                   red.Name,
			AnnualRedemptionRate:   red.AnnualRedemptionRate,
			UnitaryRedemptionValue: float64(red.UnitaryRedemptionValue),
		}
	}

	for _, inv := range investData {
		key := strings.ToLower(inv.Name)
		bond := bondsMap[key]

		bond.Name = inv.Name
		bond.Investable = true
		bond.AnnualInvestmentRate = inv.AnnualInvestmentRate
		bond.UnitaryInvestmentValue = float64(inv.UnitaryInvestmentValue)
		bond.MinimumInvestmentAmount = float64(inv.MinimumInvestmentAmount)

		bondsMap[key] = bond
	}

	bondsList := make([]ScraperBond, 0, len(bondsMap))
	for _, bond := range bondsMap {
		bondsList = append(bondsList, bond)
	}

	sc.BondsMap = bondsMap
	sc.BondsList = sc.sortByPrefix(bondsList)
	sc.UpdatedAt = time.Now().Format(time.RFC3339)
}

func (sc *ScraperCache) sortByPrefix(bonds []ScraperBond) []ScraperBond {
	order := strings.Split(viper.GetString("SORT_ORDER"), ",")

	sort.Slice(bonds, func(i, j int) bool {
		iIndex, jIndex := prefixIndex(order, bonds[i].Name), prefixIndex(order, bonds[j].Name)

		if iIndex == jIndex {
			return bonds[i].Name < bonds[j].Name
		}

		return iIndex < jIndex
	})

	return bonds
}

func prefixIndex(order []string, name string) int {
	for i, prefix := range order {
		if strings.HasPrefix(name, prefix) {
			return i
		}
	}

	return len(order)
}
