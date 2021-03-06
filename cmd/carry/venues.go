package main

import (
	"os"

	"github.com/stevenwilkin/carry/binance"
	"github.com/stevenwilkin/carry/binance_futures"
	"github.com/stevenwilkin/carry/bybit"
	"github.com/stevenwilkin/carry/deribit"

	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
)

var _deribit *deribit.Deribit
var bfRemainder int

func newBinance() *binance.Binance {
	return &binance.Binance{
		ApiKey:    os.Getenv("BINANCE_API_KEY"),
		ApiSecret: os.Getenv("BINANCE_API_SECRET")}
}

// NOTE: contract size is 100 USD
func newBinanceFutures() *binance_futures.BinanceFutures {
	return &binance_futures.BinanceFutures{
		ApiKey:    os.Getenv("BINANCE_API_KEY"),
		ApiSecret: os.Getenv("BINANCE_API_SECRET"),
		Testnet:   testnet()}
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
	}

	return _deribit
}

func bybitLimitTrader(buy bool) limitTrader {
	b := newBybit()

	return func(amount int, mt marketTrader) error {
		return b.Trade(amount, buy, buy, mt)
	}
}

func bybitMarketTrader() marketTrader {
	b := newBybit()

	return func(contracts int) {
		if err := b.MarketOrder(contracts, true, true); err != nil {
			log.Fatal(err)
		}
	}
}

func deribitLimitTrader(contract string, buy bool) limitTrader {
	d := newDeribit()

	return func(amount int, mt marketTrader) error {
		return d.Trade(contract, amount, buy, buy, mt)
	}
}

func deribitMarketTrader(contract string) marketTrader {
	d := newDeribit()

	return func(amount int) {
		if err := d.MarketOrder(contract, amount, true, true); err != nil {
			log.Fatal(err)
		}
	}
}

func binanceMarketTrader(buy bool) marketTrader {
	b := newBinance()

	return func(usdt int) {
		if err := b.MarketOrder(float64(usdt), buy); err != nil {
			log.Fatal(err)
		}
	}
}

func binanceFuturesLimitTrader(buy bool) limitTrader {
	b := newBinanceFutures()

	return func(amount int, mt marketTrader) error {
		return b.Trade(amount/100, buy, buy, mt)
	}
}

func binanceFuturesMarketTrader() marketTrader {
	b := newBinanceFutures()

	return func(amount int) {
		contracts := (amount + bfRemainder) / 100
		bfRemainder = (amount + bfRemainder) % 100

		if contracts == 0 {
			log.WithFields(log.Fields{
				"venue":     "binance_f",
				"remainder": bfRemainder,
			}).Debug("Amount below contract size")
			return
		}

		if err := b.MarketOrder(contracts, true, true); err != nil {
			log.Fatal(err)
		}
	}
}
