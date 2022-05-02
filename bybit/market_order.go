package bybit

import (
	"strconv"

	log "github.com/sirupsen/logrus"
)

func (b *Bybit) MarketOrder(contracts int, buy, reduce bool) string {
	log.WithFields(log.Fields{
		"venue":     "bybit",
		"contracts": contracts,
		"buy":       buy,
	}).Debug("Placing market order")

	params := map[string]interface{}{
		"symbol":        "BTCUSD",
		"order_type":    "Market",
		"qty":           strconv.Itoa(contracts),
		"time_in_force": "GoodTillCancel"}

	if buy {
		params["side"] = "Buy"
	} else {
		params["side"] = "Sell"
	}

	if reduce {
		params["reduce_only"] = true
	}

	var result orderResponse
	if err := b.post("/v2/private/order/create", params, &result); err != nil {
		log.Error(err.Error())
	}

	return result.Result.OrderId
}
