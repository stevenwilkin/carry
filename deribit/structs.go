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

type orderResponse struct {
	Result struct {
		Order struct {
			OrderId string `json:"order_id"`
		} `json:"order"`
	} `json:"result"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

type requestMessage struct {
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
}

type quoteMessage struct {
	Method string `json:"method"`
	Params struct {
		Data struct {
			BestBidPrice float64 `json:"best_bid_price"`
			BestAskPrice float64 `json:"best_ask_price"`
		} `json:"data"`
	} `json:"params"`
}
