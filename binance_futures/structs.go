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
	OrderId int64 `json:"orderId"`
}
