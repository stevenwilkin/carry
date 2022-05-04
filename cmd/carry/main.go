package main

import (
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

		if contract == "BTCUSD" {
			lt = bybitLimitTrader(buyContracts)
		} else {
			lt = deribitLimitTrader(contract, buyContracts)
		}
	} else if action == "roll" {
		if rollToContract == "BTCUSD" { // deribit -> bybit
			lt = bybitLimitTrader(false)
			mt = deribitMarketTrader(contract)
		} else {
			if contract == "BTCUSD" { // bybit -> deribit
				mt = bybitMarketTrader()
			} else { // deribit -> deribit
				mt = deribitMarketTrader(contract)
			}

			lt = deribitLimitTrader(rollToContract, false)
		}
	}

	for i := 0; i < rounds; i++ {
		log.WithField("n", i+1).Info("Round")

		if err := lt(usd, mt); err != nil {
			log.Fatal(err)
		}
	}

	log.Info("Done")
}
