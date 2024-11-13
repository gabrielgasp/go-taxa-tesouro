package model

import "strings"

type ScrapperBond struct {
	Name                    string  `json:"name"`
	AnnualInvestmentRate    float64 `json:"annual_investment_rate"`
	UnitaryInvestmentValue  float64 `json:"unitary_investment_value"`
	MinimumInvestmentAmount float64 `json:"minimum_investment_amount"`
	AnnualRedemptionRate    float64 `json:"annual_redemption_rate"`
	UnitaryRedemptionValue  float64 `json:"unitary_redemption_value"`
	MinimumRedemptionValue  float64 `json:"minimum_redemption_value"`
}

type ScrapperCache struct {
	BondsList []ScrapperBond
	BondsMap  map[string]ScrapperBond
	UpdatedAt string
}

func (sc *ScrapperCache) Save(data TesouroData) {
	sc.BondsList = make([]ScrapperBond, len(data.TreasureBondsList))
	sc.BondsMap = make(map[string]ScrapperBond)

	for i, tb := range data.TreasureBondsList {
		scrapperBond := ScrapperBond{
			Name:                    tb.TreasureBond.Nm,
			AnnualInvestmentRate:    tb.TreasureBond.AnulInvstmtRate,
			UnitaryInvestmentValue:  tb.TreasureBond.UntrInvstmtVal,
			MinimumInvestmentAmount: tb.TreasureBond.MinInvstmtAmt,
			AnnualRedemptionRate:    tb.TreasureBond.AnulRedRate,
			UnitaryRedemptionValue:  tb.TreasureBond.UntrRedVal,
			MinimumRedemptionValue:  tb.TreasureBond.MinRedVal,
		}

		sc.BondsList[i] = scrapperBond
		sc.BondsMap[strings.ToLower(scrapperBond.Name)] = scrapperBond
	}
	sc.UpdatedAt = data.BusinessStatus.DtTm
}
