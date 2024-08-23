package model

type TesouroResponse struct {
	Response struct {
		TreasureBondsList []struct {
			TreasureBond TreasureBond `json:"TrsrBd"`
		} `json:"TrsrBdTradgList"`
		BusinessStatus BusinessStatus `json:"BizSts"`
	} `json:"response"`
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
