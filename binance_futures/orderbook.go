package binance_futures

import (
	"strconv"

	log "github.com/sirupsen/logrus"
)

func (b *BinanceFutures) SubscribeToOrderBook() error {
	c, err := b.subscribe("btcusd_perp@bookTicker")
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	go func() {
		defer c.Close()
		var bookTicker bookTickerMessage

		for {
			err := c.ReadJSON(&bookTicker)
			if err != nil {
				log.Error(err.Error())
				return
			}

			b.Bid, _ = strconv.ParseFloat(bookTicker.BidPrice, 64)
			b.Ask, _ = strconv.ParseFloat(bookTicker.AskPrice, 64)
		}
	}()

	for b.Bid == 0 || b.Ask == 0 {
	}

	return nil
}
