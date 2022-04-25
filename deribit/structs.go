package deribit

type authResponse struct {
	Result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	} `json:"result"`
}

type positionsResponse struct {
	Result []struct {
		InstrumentName string  `json:"instrument_name"`
		Size           float64 `json:"size"`
	} `json:"result"`
}
