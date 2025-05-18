package deribit

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

func (d *Deribit) subscribeToOrders(instrument string) chan float64 {
	completed := make(chan float64)
	ordersChannel := fmt.Sprintf("user.orders.%s.raw", instrument)

	c, err := d.subscribe([]string{ordersChannel})
	if err != nil {
		return completed
	}

	go func() {
		var m orderMessage
		defer c.Close()

		for {
			if err = c.ReadJSON(&m); err != nil {
				log.Error(err.Error())
				return
			}

			if m.Method != "subscription" {
				continue
			}

			switch m.Params.Data.OrderState {
			case "open":
				log.WithFields(log.Fields{
					"venue":    "deribit",
					"order":    m.Params.Data.OrderId,
					"quantity": m.Params.Data.FilledAmount,
				}).Debug("Order open")
				completed <- m.Params.Data.FilledAmount
			case "cancelled":
				log.WithFields(log.Fields{
					"venue": "deribit",
					"order": m.Params.Data.OrderId,
				}).Debug("Order cancelled")
				return
			case "filled":
				log.WithFields(log.Fields{
					"venue":    "deribit",
					"order":    m.Params.Data.OrderId,
					"quantity": m.Params.Data.FilledAmount,
				}).Debug("Order filled")
				completed <- m.Params.Data.FilledAmount
				close(completed)
				return
			}
		}
	}()

	return completed
}

func (d *Deribit) makeBestPrice(buy bool) func() float64 {
	if buy {
		return func() float64 {
			return d.Bid
		}
	} else {
		return func() float64 {
			return d.Ask
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

func (d *Deribit) Trade(instrument string, contracts int, buy, reduce bool, cb func(int)) error {
	log.WithFields(log.Fields{
		"venue":      "deribit",
		"instrument": instrument,
		"contracts":  contracts,
		"buy":        buy,
		"reduce":     reduce,
	}).Info("Trade")

	if err := d.SubscribeToOrderBook(instrument); err != nil {
		log.Fatal(err)
	}

	var price, newBestPrice, totalCompleted float64
	var orderId string
	var err error
	bestPrice := d.makeBestPrice(buy)
	ch := d.subscribeToOrders(instrument)
	ticker := time.NewTicker(10 * time.Millisecond)

	for {
		select {
		case completed, ok := <-ch:
			if !ok {
				log.WithField("venue", "deribit").Info("Trade completed")
				return nil
			}
			if completed > totalCompleted {
				delta := completed - totalCompleted
				log.WithFields(log.Fields{
					"venue":  "deribit",
					"amount": delta,
				}).Debug("Callback")
				cb(int(delta))
				totalCompleted = completed
			}
		case <-ticker.C:
			if orderId == "" {
				price = bestPrice()
				orderId, err = d.LimitOrder(instrument, contracts, price, buy, reduce)
				if err != nil {
					return err
				}
			} else {
				newBestPrice = bestPrice()
				if canImprove(price, newBestPrice, buy) {
					price = newBestPrice
					if err = d.EditOrder(orderId, contracts, price, reduce); err != nil {
						log.Warn(err)
					}
				}
			}
		}
	}

	return nil
}
