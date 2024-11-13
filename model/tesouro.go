package model

type TesouroResponse struct {
	Data TesouroData `json:"response"`
}

type TesouroData struct {
	TreasureBondsList []TreasureBondsList `json:"TrsrBdTradgList"`
	BusinessStatus    BusinessStatus      `json:"BizSts"`
}

type TreasureBondsList struct {
	TreasureBond TreasureBond `json:"TrsrBd"`
}

type TreasureBond struct {
	Nm              string  `json:"nm"`
	AnulInvstmtRate float64 `json:"anulInvstmtRate"`
	UntrInvstmtVal  float64 `json:"untrInvstmtVal"`
	MinInvstmtAmt   float64 `json:"minInvstmtAmt"`
	AnulRedRate     float64 `json:"anulRedRate"`
	UntrRedVal      float64 `json:"untrRedVal"`
	MinRedVal       float64 `json:"minRedVal"`
}

type BusinessStatus struct {
	DtTm string `json:"dtTm"`
}
