package bybit

import (
	"errors"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func (b *Bybit) LimitOrder(contracts int, price float64, buy, reduce bool) (string, error) {
	log.WithFields(log.Fields{
		"venue":     "bybit",
		"contracts": contracts,
		"price":     price,
		"buy":       buy,
	}).Debug("Placing order")

	params := map[string]interface{}{
		"category":    "inverse",
		"symbol":      "BTCUSD",
		"orderType":   "Limit",
		"qty":         strconv.Itoa(contracts),
		"price":       strconv.FormatFloat(price, 'f', 2, 64),
		"timeInForce": "PostOnly"}

	if buy {
		params["side"] = "Buy"
	} else {
		params["side"] = "Sell"
	}

	if reduce {
		params["reduceOnly"] = true
	}

	var result orderResponse
	if err := b.post("/v5/order/create", params, &result); err != nil {
		return "", err
	}

	if result.RetCode != 0 {
		return "", errors.New(result.RetMsg)
	}

	return result.Result.OrderId, nil
}

func (b *Bybit) EditOrder(id string, price float64) error {
	log.WithFields(log.Fields{
		"venue": "bybit",
		"order": id,
		"price": price,
	}).Debug("Updating order")

	params := map[string]interface{}{
		"category": "inverse",
		"orderId":  id,
		"symbol":   "BTCUSD",
		"price":    strconv.FormatFloat(price, 'f', 2, 64)}

	var result orderResponse
	if err := b.post("/v5/order/amend", params, &result); err != nil {
		return err
	}

	if result.RetCode != 0 {
		return errors.New(result.RetMsg)
	}

	return nil
}

func (b *Bybit) CancelOrders() {
	log.WithField("venue", "bybit").Info("Cancelling orders")

	params := map[string]interface{}{
		"category": "inverse",
		"symbol":   "BTCUSD"}

	var result orderResponse
	if err := b.post("/v5/order/cancel-all", params, &result); err != nil {
		log.Error(err)
		return
	}

	if result.RetCode != 0 {
		log.Error(result.RetMsg)
	}
}
