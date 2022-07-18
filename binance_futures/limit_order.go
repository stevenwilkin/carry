package binance_futures

import (
	"encoding/json"
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func (b *BinanceFutures) LimitOrder(contracts int, price float64, buy, reduce bool) (int64, error) {
	log.WithFields(log.Fields{
		"venue":     "binance_f",
		"contracts": contracts,
		"price":     price,
		"buy":       buy,
	}).Debug("Placing order")

	params := url.Values{
		"symbol":      {"BTCUSD_PERP"},
		"type":        {"LIMIT"},
		"quantity":    {strconv.Itoa(contracts)},
		"price":       {strconv.FormatFloat(price, 'f', 2, 64)},
		"timeInForce": {"GTX"}}

	if buy {
		params.Add("side", "BUY")
	} else {
		params.Add("side", "SELL")
	}

	if reduce {
		params.Add("reduceOnly", "true")
	}

	body, err := b.doRequest("POST", "/dapi/v1/order", params, true)
	if err != nil {
		return 0, err
	}

	var result orderResponse
	json.Unmarshal(body, &result)

	return result.OrderId, nil
}
