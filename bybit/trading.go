package bybit

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func (b *Bybit) subscribeToOrders() (chan int, chan bool) {
	ch := make(chan int)
	chCancels := make(chan bool)
	orderTopic := "order"

	c, err := b.subscribe([]string{orderTopic})
	if err != nil {
		log.Error(err.Error())
		return ch, chCancels
	}

	go func() {
		defer c.Close()
		var orders orderTopicData

		for {
			if err := c.ReadJSON(&orders); err != nil {
				log.Error(err)
				return
			}

			if orders.Topic != orderTopic {
				continue
			}

			order := orders.Data[0]

			switch order.OrderStatus {
			case "PartiallyFilled":
				log.WithFields(log.Fields{
					"venue":        "bybit",
					"order":        order.OrderId,
					"quantity":     order.Qty,
					"cum_quantity": order.CumExecQty,
				}).Debug("Fill")
				ch <- order.CumExecQty
			case "Filled":
				log.WithFields(log.Fields{
					"venue":        "bybit",
					"order":        order.OrderId,
					"quantity":     order.Qty,
					"cum_quantity": order.CumExecQty,
				}).Debug("Order filled")
				ch <- order.CumExecQty
				close(ch)
				return
			case "Cancelled":
				log.WithFields(log.Fields{
					"venue":        "bybit",
					"order":        order.OrderId,
					"quantity":     order.Qty,
					"cum_quantity": order.CumExecQty,
				}).Debug("Order cancelled")
				chCancels <- true
			}
		}
	}()

	return ch, chCancels
}

func (b *Bybit) makeBestPrice(buy bool) func() float64 {
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

func (b *Bybit) canImprove(price, bestPrice float64, buy bool) bool {
	if buy {
		return price < bestPrice
	} else {
		return price > bestPrice
	}
}

func (b *Bybit) Trade(contracts int, buy, reduce bool, cb func(int)) {
	log.WithFields(log.Fields{
		"venue":     "bybit",
		"contracts": contracts,
		"buy":       buy,
		"reduce":    reduce,
	}).Info("Trade")

	b.o.Do(func() {
		if err := b.SubscribeToOrderBook(); err != nil {
			log.Fatal(err)
		}
	})

	var orderId string
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
				log.WithField("venue", "bybit").Info("Trade completed")
				return
			}
			if completed > totalCompleted {
				delta := completed - totalCompleted
				log.WithFields(log.Fields{
					"venue":  "bybit",
					"amount": delta,
				}).Debug("Callback")
				cb(delta)
				totalCompleted = completed
			}
		case <-chCancels:
			orderId = ""
		case <-ticker.C:
			if orderId == "" {
				price = bestPrice()
				orderId, err = b.LimitOrder(contracts-totalCompleted, price, buy, reduce)
				if err != nil {
					log.Error(err)
					return
				}
			} else {
				newBestPrice = bestPrice()
				if b.canImprove(price, newBestPrice, buy) {
					price = newBestPrice
					if err = b.EditOrder(orderId, newBestPrice); err != nil {
						log.Error(err)
						return
					}
				}
			}
		}
	}
}
