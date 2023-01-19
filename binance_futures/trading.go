package binance_futures

import (
	"encoding/json"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

func (b *BinanceFutures) listenKey() (string, error) {
	body, err := b.doRequest("POST", "/dapi/v1/listenKey", nil, false)
	if err != nil {
		return "", err
	}

	var response listenKeyResponse
	json.Unmarshal(body, &response)

	return response.ListenKey, nil
}

func (b *BinanceFutures) subscribeToOrders() (chan int, chan bool) {
	ch := make(chan int)
	chCancels := make(chan bool)

	key, err := b.listenKey()
	if err != nil {
		log.Error(err.Error())
		return ch, chCancels
	}

	c, err := b.subscribe(key)
	if err != nil {
		log.Error(err.Error())
		return ch, chCancels
	}

	go func() {
		defer c.Close()
		var m userDataMessage
		var orderId, fillQty, cumFillQty int

		for {
			if err := c.ReadJSON(&m); err != nil {
				log.Error(err.Error())
				continue
			}

			if m.EventType != "ORDER_TRADE_UPDATE" {
				continue
			}

			orderId = m.Order.OrderId
			fillQty, _ = strconv.Atoi(m.Order.FillQty)
			cumFillQty, _ = strconv.Atoi(m.Order.CumFillQty)

			switch m.Order.OrderStatus {
			case "PARTIALLY_FILLED":
				log.WithFields(log.Fields{
					"venue":        "binance_f",
					"order":        orderId,
					"quantity":     fillQty,
					"cum_quantity": cumFillQty,
				}).Debug("Fill")
				ch <- cumFillQty
			case "FILLED":
				log.WithFields(log.Fields{
					"venue":        "binance_f",
					"order":        orderId,
					"quantity":     fillQty,
					"cum_quantity": cumFillQty,
				}).Debug("Order filled")
				ch <- cumFillQty
				close(ch)
				return
			case "EXPIRED", "CANCELED":
				log.WithFields(log.Fields{
					"venue":        "binance_f",
					"order":        orderId,
					"cum_quantity": cumFillQty,
				}).Debug("Order cancelled")
				chCancels <- true
			}
		}
	}()

	return ch, chCancels
}

func (b *BinanceFutures) makeBestPrice(buy bool) func() float64 {
	if buy {
		return func() float64 {
			return b.Bid
		}
	} else {
		return func() float64 {
			return b.Ask
		}
	}
}

func canImprove(price, bestPrice float64, buy bool) bool {
	if buy {
		return price < bestPrice
	} else {
		return price > bestPrice
	}
}

func (b *BinanceFutures) Trade(contracts int, buy, reduce bool, cb func(int)) error {
	log.WithFields(log.Fields{
		"venue":     "binance_f",
		"contracts": contracts,
		"buy":       buy,
		"reduce":    reduce,
	}).Info("Trade")

	b.o.Do(func() {
		if err := b.SubscribeToOrderBook(); err != nil {
			log.Fatal(err)
		}
	})

	var orderId int
	var price, newBestPrice float64
	var totalCompleted int
	var err error

	bestPrice := b.makeBestPrice(buy)
	ch, chCancels := b.subscribeToOrders()
	ticker := time.NewTicker(10 * time.Millisecond)

	for {
		select {
		case completed, ok := <-ch:
			if !ok {
				log.WithField("venue", "binance_f").Info("Trade completed")
				return nil
			}
			if completed > totalCompleted {
				delta := (completed - totalCompleted) * 100
				log.WithFields(log.Fields{
					"venue":  "binance_f",
					"amount": delta,
				}).Debug("Callback")
				go cb(delta)
				totalCompleted = completed
			}
		case <-chCancels:
			orderId = 0
		case <-ticker.C:
			if orderId == 0 {
				price = bestPrice()
				orderId, err = b.LimitOrder(contracts-totalCompleted, price, buy, reduce)
				if err != nil {
					return err
				}
			} else {
				newBestPrice = bestPrice()
				if canImprove(price, newBestPrice, buy) {
					price = newBestPrice
					if err = b.EditOrder(orderId, newBestPrice, buy); err != nil {
						log.Warn(err)
					}
				}
			}
		}
	}
}
