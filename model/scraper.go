package model

import "strings"

type ScraperBond struct {
	Name                    string  `json:"name"`
	AnnualInvestmentRate    float64 `json:"annual_investment_rate"`
	UnitaryInvestmentValue  float64 `json:"unitary_investment_value"`
	MinimumInvestmentAmount float64 `json:"minimum_investment_amount"`
	AnnualRedemptionRate    float64 `json:"annual_redemption_rate"`
	UnitaryRedemptionValue  float64 `json:"unitary_redemption_value"`
	MinimumRedemptionValue  float64 `json:"minimum_redemption_value"`
}

type ScraperCache struct {
	BondsList []ScraperBond
	BondsMap  map[string]ScraperBond
	UpdatedAt string
}

func (sc *ScraperCache) Save(data TesouroData) {
	sc.BondsList = make([]ScraperBond, len(data.TreasureBondsList))
	sc.BondsMap = make(map[string]ScraperBond)

	for i, tb := range data.TreasureBondsList {
		scraperBond := ScraperBond{
			Name:                    tb.TreasureBond.Nm,
			AnnualInvestmentRate:    tb.TreasureBond.AnulInvstmtRate,
			UnitaryInvestmentValue:  tb.TreasureBond.UntrInvstmtVal,
			MinimumInvestmentAmount: tb.TreasureBond.MinInvstmtAmt,
			AnnualRedemptionRate:    tb.TreasureBond.AnulRedRate,
			UnitaryRedemptionValue:  tb.TreasureBond.UntrRedVal,
			MinimumRedemptionValue:  tb.TreasureBond.MinRedVal,
		}

		sc.BondsList[i] = scraperBond
		sc.BondsMap[strings.ToLower(scraperBond.Name)] = scraperBond
	}
	sc.UpdatedAt = data.BusinessStatus.DtTm
}
