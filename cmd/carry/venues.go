package main

import (
	"github.com/stevenwilkin/carry/binance"
	"github.com/stevenwilkin/carry/bybit"
	"github.com/stevenwilkin/carry/deribit"

	_ "github.com/joho/godotenv/autoload"
)

func bybitLimitTrader(buy bool) (limitTrader, orderCanceler) {
	b := bybit.NewBybitFromEnv()

	return func(amount int, mt cb) error {
			return b.Trade(amount, buy, buy, mt)
		}, func() {
			b.CancelOrders()
		}
}

func bybitMarketTrader(buy bool) marketTrader {
	b := bybit.NewBybitFromEnv()

	return func(contracts int) error {
		return b.MarketOrder(contracts, buy, buy)
	}
}

func deribitLimitTrader(contract string, buy bool) (limitTrader, orderCanceler) {
	d := deribit.NewDeribitFromEnv()

	return func(amount int, mt cb) error {
			return d.Trade(contract, amount, buy, buy, mt)
		},
		func() {
			d.CancelOrders(contract)
		}
}

func deribitMarketTrader(contract string, buy bool) marketTrader {
	d := deribit.NewDeribitFromEnv()

	return func(amount int) error {
		return d.MarketOrder(contract, amount, buy, buy)
	}
}

func binanceMarketTrader(buy bool) marketTrader {
	b := binance.NewBinanceFromEnv()

	return func(amount int) error {
		return b.MarketOrder(float64(amount), buy)
	}
}
