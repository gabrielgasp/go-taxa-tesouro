package model

type ListAllBondsResponse struct {
	Bonds     []ScraperBond `json:"bonds"`
	UpdatedAt string        `json:"updated_at"`
}

type GetBondByNameResponse struct {
	Bond      ScraperBond `json:"bond"`
	UpdatedAt string      `json:"updated_at"`
}
