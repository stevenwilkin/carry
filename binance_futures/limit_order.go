package binance_futures

import (
	"encoding/json"
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func (b *BinanceFutures) LimitOrder(contracts int, price float64, buy, reduce bool) (int, error) {
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

func (b *BinanceFutures) EditOrder(id int, price float64, buy bool) error {
	log.WithFields(log.Fields{
		"venue": "binance_f",
		"order": id,
		"price": price,
	}).Debug("Updating order")

	params := url.Values{
		"orderId": {strconv.Itoa(id)},
		"symbol":  {"BTCUSD_PERP"},
		"price":   {strconv.FormatFloat(price, 'f', 2, 64)}}

	if buy {
		params.Add("side", "BUY")
	} else {
		params.Add("side", "SELL")
	}

	if _, err := b.doRequest("PUT", "/dapi/v1/order", params, true); err != nil {
		return err
	}

	return nil
}
