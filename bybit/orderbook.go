package bybit

import (
	"encoding/json"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func highest(orders map[int64]float64) float64 {
	var result float64

	for _, x := range orders {
		if x > result {
			result = x
		}
	}

	return result
}

func lowest(orders map[int64]float64) float64 {
	var result float64

	for _, x := range orders {
		if result == 0.0 {
			result = x
		} else if x < result {
			result = x
		}
	}

	return result
}

func (b *Bybit) SubscribeToOrderBook() error {
	bids := map[int64]float64{}
	asks := map[int64]float64{}

	orderBookTopic := "orderBookL2_25.BTCUSD"

	c, err := b.subscribe([]string{orderBookTopic})
	if err != nil {
		log.Error(err.Error())
		return err
	}

	go func() {
		defer c.Close()
		var response wsResponse

		for {
			if err := c.ReadJSON(&response); err != nil {
				log.Error(err)
				return
			}

			if response.Topic != orderBookTopic {
				continue
			}

			switch response.Type {
			case "snapshot":
				var snapshot snapshotData
				json.Unmarshal(response.Data, &snapshot)

				for _, order := range snapshot {
					p, _ := strconv.ParseFloat(order.Price, 64)

					if order.Side == "Buy" {
						bids[order.Id] = p
					} else {
						asks[order.Id] = p
					}
				}
			case "delta":
				var updates updateData
				json.Unmarshal(response.Data, &updates)

				for _, order := range updates.Delete {
					if order.Side == "Buy" {
						delete(bids, order.Id)
					} else {
						delete(asks, order.Id)
					}
				}

				for _, order := range updates.Insert {
					p, _ := strconv.ParseFloat(order.Price, 64)

					if order.Side == "Buy" {
						bids[order.Id] = p
					} else {
						asks[order.Id] = p
					}
				}
			}

			b.Bid = highest(bids)
			b.Ask = lowest(asks)
		}
	}()

	for b.Bid == 0 || b.Ask == 0 {
	}

	return nil
}
