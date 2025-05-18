package main

import (
	"github.com/stevenwilkin/carry/binance"
	"github.com/stevenwilkin/carry/bybit"
	"github.com/stevenwilkin/carry/deribit"

	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
)

func bybitLimitTrader(buy bool) (limitTrader, orderCanceler) {
	b := bybit.NewBybitFromEnv()

	return func(amount int, mt marketTrader) error {
			return b.Trade(amount, buy, buy, mt)
		}, func() {
			b.CancelOrders()
		}
}

func bybitMarketTrader(buy bool) marketTrader {
	b := bybit.NewBybitFromEnv()

	return func(contracts int) {
		if err := b.MarketOrder(contracts, buy, buy); err != nil {
			log.Error(err)
		}
	}
}

func deribitLimitTrader(contract string, buy bool) (limitTrader, orderCanceler) {
	d := deribit.NewDeribitFromEnv()

	return func(amount int, mt marketTrader) error {
			return d.Trade(contract, amount, buy, buy, mt)
		},
		func() {
			d.CancelOrders(contract)
		}
}

func deribitMarketTrader(contract string, buy bool) marketTrader {
	var dRemainder int
	d := deribit.NewDeribitFromEnv()

	return func(amount int) {
		// must be multiples of 10
		contracts := ((amount + dRemainder) / 10) * 10
		dRemainder = (amount + dRemainder) % 10

		log.WithFields(log.Fields{
			"venue":     "deribit",
			"amount":    amount,
			"contracts": contracts,
			"remainder": dRemainder,
		}).Debug("Market order")

		if contracts == 0 {
			log.WithFields(log.Fields{
				"venue":     "deribit",
				"remainder": dRemainder,
			}).Info("Skipping market order")
			return
		}
		if err := d.MarketOrder(contract, contracts, buy, buy); err != nil {
			log.Error(err)
		}
	}
}

func binanceMarketTrader(buy bool) marketTrader {
	var bRemainder int
	b := binance.NewBinanceFromEnv()

	return func(usdt int) {
		// must be greater than 10
		amount := usdt + bRemainder

		if amount < 10 {
			bRemainder = amount
			amount = 0
		} else {
			bRemainder = 0
		}

		log.WithFields(log.Fields{
			"venue":     "binance",
			"usdt":      usdt,
			"amount":    amount,
			"remainder": bRemainder,
		}).Debug("Market order")

		if amount == 0 {
			log.WithFields(log.Fields{
				"venue":     "binance",
				"remainder": bRemainder,
			}).Info("Skipping market order")
			return
		}

		if err := b.MarketOrder(float64(amount), buy); err != nil {
			log.Error(err)
		}
	}
}
