package main

import (
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type marketTrader func(amount int)
type limitTrader func(amount int, mt marketTrader) error

var (
	params                           []string
	action, contract, rollToContract string
	usd, rounds                      int
	lt                               limitTrader
	mt                               marketTrader
)

func main() {
	initParams()

	log.WithFields(log.Fields{
		"action":           action,
		"contract":         contract,
		"roll_to_contract": rollToContract,
		"usd":              usd,
		"rounds":           rounds,
	}).Debug("Params")

	if action == "up" || action == "down" {
		buyContracts := action == "down"
		mt = binanceMarketTrader(!buyContracts)

		switch contract {
		case "BTCUSD":
			lt = bybitLimitTrader(buyContracts)
		case "BTCUSD_PERP":
			lt = binanceFuturesLimitTrader(buyContracts)
		default:
			lt = deribitLimitTrader(contract, buyContracts)
		}
	} else if action == "roll" {
		switch contract {
		case "BTCUSD":
			mt = bybitMarketTrader(true)
		case "BTCUSD_PERP":
			mt = binanceFuturesMarketTrader()
		default:
			mt = deribitMarketTrader(contract, true)
		}

		switch rollToContract {
		case "BTCUSD":
			lt = bybitLimitTrader(false)
		case "BTCUSD_PERP":
			lt = binanceFuturesLimitTrader(false)
		default:
			lt = deribitLimitTrader(rollToContract, false)
		}
	} else if action == "rollx" {
		// roll from less liquid to more liquid contracts
		if !strings.HasPrefix(contract, "BTC-") {
			log.Fatal("Can only roll from Deribit")
		}

		lt = deribitLimitTrader(contract, true)

		if rollToContract == "BTCUSD" {
			mt = bybitMarketTrader(false)
		} else if rollToContract == "BTCUSD_PERP" {
			log.Fatal("Can not roll to Binance Futures")
		} else {
			mt = deribitMarketTrader(rollToContract, false)
		}
	}

	for i := 0; i < rounds; i++ {
		log.WithField("n", i+1).Info("Round")

		if err := lt(usd, mt); err != nil {
			log.Fatal(err)
		}
	}

	log.Debug("Sleep")
	time.Sleep(3 * time.Second)

	log.Info("Done")
}
