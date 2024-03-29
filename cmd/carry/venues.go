package main

import (
	"os"

	"github.com/stevenwilkin/carry/binance"
	"github.com/stevenwilkin/carry/bybit"
	"github.com/stevenwilkin/carry/deribit"

	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
)

var _deribit *deribit.Deribit

func newBinance() *binance.Binance {
	return &binance.Binance{
		ApiKey:    os.Getenv("BINANCE_API_KEY"),
		ApiSecret: os.Getenv("BINANCE_API_SECRET")}
}

func newBybit() *bybit.Bybit {
	return &bybit.Bybit{
		ApiKey:    os.Getenv("BYBIT_API_KEY"),
		ApiSecret: os.Getenv("BYBIT_API_SECRET"),
		Testnet:   testnet()}
}

func newDeribit() *deribit.Deribit {
	if _deribit == nil {
		_deribit = &deribit.Deribit{
			ApiId:     os.Getenv("DERIBIT_API_ID"),
			ApiSecret: os.Getenv("DERIBIT_API_SECRET"),
			Test:      testnet()}

		if _, err := _deribit.AccessToken(); err != nil {
			log.Fatal(err)
		}
	}

	return _deribit
}

func bybitLimitTrader(buy bool) (limitTrader, orderCanceler) {
	b := newBybit()

	return func(amount int, mt marketTrader) error {
			return b.Trade(amount, buy, buy, mt)
		}, func() {
			b.CancelOrders()
		}
}

func bybitMarketTrader(buy bool) marketTrader {
	b := newBybit()

	return func(contracts int) {
		if err := b.MarketOrder(contracts, buy, buy); err != nil {
			log.Error(err)
		}
	}
}

func deribitLimitTrader(contract string, buy bool) (limitTrader, orderCanceler) {
	d := newDeribit()

	return func(amount int, mt marketTrader) error {
			return d.Trade(contract, amount, buy, buy, mt)
		},
		func() {
			d.CancelOrders(contract)
		}
}

func deribitMarketTrader(contract string, buy bool) marketTrader {
	var dRemainder int
	d := newDeribit()

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
	b := newBinance()

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
