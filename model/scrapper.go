package model

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
