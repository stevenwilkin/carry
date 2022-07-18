package binance_futures

type bookTickerMessage struct {
	BidPrice string `json:"b"`
	BidQty   string `json:"B"`
	AskPrice string `json:"a"`
	AskQty   string `json:"A"`
}

type errorResponse struct {
	Msg string `json:"msg"`
}

type orderResponse struct {
	OrderId int `json:"orderId"`
}

type listenKeyResponse struct {
	ListenKey string `json:"listenKey"`
}

type userDataMessage struct {
	EventType string `json:"e"`
	EventTime int64  `json:"E"`
	Order     struct {
		ExecutionType string `json:"x"`
		OrderStatus   string `json:"X"`
		OrderId       int    `json:"i"`
		FillQty       string `json:"l"`
		FillPrice     string `json:"L"`
		CumFillQty    string `json:"z"`
	} `json:"o"`
}
