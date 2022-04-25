package bybit

type positionResponse struct {
	Result struct {
		Size int `json:"size"`
	} `json:"result"`
}
