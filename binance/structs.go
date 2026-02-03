package binance

import "strconv"

type assetBalance struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

func (ab *assetBalance) Total() float64 {
	free, _ := strconv.ParseFloat(ab.Free, 64)
	locked, _ := strconv.ParseFloat(ab.Locked, 64)
	return free + locked
}

type accountResponse struct {
	Balances []assetBalance `json:"balances"`
}

type errorResponse struct {
	Msg string `json:"msg"`
}

type addressResponse struct {
	Address string `json:"address"`
}

type balance struct {
	Asset   string
	Balance float64
	Value   float64
}

type priceResponse struct {
	Price string `json:"price"`
}

func (p *priceResponse) PriceFloat() float64 {
	price, _ := strconv.ParseFloat(p.Price, 64)
	return price
}
