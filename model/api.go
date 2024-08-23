package model

type ListAllBondsResponse struct {
	Bonds     []ScrapperBond `json:"bonds"`
	UpdatedAt string         `json:"updated_at"`
}

type GetBondByNameResponse struct {
	Bond      ScrapperBond `json:"bond"`
	UpdatedAt string       `json:"updated_at"`
}
