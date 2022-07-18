package binance_futures

import (
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func (b *BinanceFutures) MarketOrder(contracts int, buy, reduce bool) error {
	log.WithFields(log.Fields{
		"venue":     "binance_f",
		"contracts": contracts,
		"buy":       buy,
	}).Debug("Placing market order")

	params := url.Values{
		"symbol":   {"BTCUSD_PERP"},
		"type":     {"MARKET"},
		"quantity": {strconv.Itoa(contracts)}}

	if buy {
		params.Add("side", "BUY")
	} else {
		params.Add("side", "SELL")
	}

	if reduce {
		params.Add("reduceOnly", "true")
	}

	_, err := b.doRequest("POST", "/dapi/v1/order", params, true)
	if err != nil {
		return err
	}

	return nil
}
