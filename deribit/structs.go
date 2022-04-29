package deribit

type authResponse struct {
	Result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	} `json:"result"`
}

type Position struct {
	InstrumentName string  `json:"instrument_name"`
	Size           float64 `json:"size"`
}

type positionsResponse struct {
	Result []Position `json:"result"`
}
