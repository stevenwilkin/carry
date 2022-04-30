package deribit

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func (d *Deribit) SubscribeToOrderBook(instrument string) error {
	quoteChannel := fmt.Sprintf("quote.%s", instrument)

	c, err := d.subscribe([]string{quoteChannel})
	if err != nil {
		log.Error(err.Error())
		return err
	}

	go func() {
		var qm quoteMessage
		defer c.Close()

		for {
			if err = c.ReadJSON(&qm); err != nil {
				log.Error(err.Error())
				return
			}

			if qm.Method != "subscription" {
				continue
			}

			d.Bid = qm.Params.Data.BestBidPrice
			d.Ask = qm.Params.Data.BestAskPrice
		}
	}()

	for d.Bid == 0 || d.Ask == 0 {
	}

	return nil
}
