package main

import (
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"
)

type marketTrader func(amount int)
type limitTrader func(amount int, mt marketTrader) error
type orderCanceler func()

var (
	params                           []string
	action, contract, rollToContract string
	usd                              int
	lt                               limitTrader
	mt                               marketTrader
	oc                               orderCanceler
	wg                               sync.WaitGroup
)

func trapSigInt() {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-c
		if oc != nil {
			oc()
		}
		os.Exit(0)
	}()
}

func processCb(cb marketTrader) marketTrader {
	return func(amount int) {
		wg.Add(1)
		go func() {
			cb(amount)
			wg.Done()
		}()
	}
}

func main() {
	initParams()
	trapSigInt()

	log.WithFields(log.Fields{
		"action":           action,
		"contract":         contract,
		"roll_to_contract": rollToContract,
		"usd":              usd,
	}).Debug("Params")

	if action == "up" || action == "down" {
		buyContracts := action == "down"
		mt = binanceMarketTrader(!buyContracts)

		if contract == "BTCUSD" {
			lt, oc = bybitLimitTrader(buyContracts)
		} else {
			lt, oc = deribitLimitTrader(contract, buyContracts)
		}
	} else if action == "roll" {
		if contract == "BTCUSD" {
			mt = bybitMarketTrader(true)
		} else {
			mt = deribitMarketTrader(contract, true)
		}

		if rollToContract == "BTCUSD" {
			lt, oc = bybitLimitTrader(false)
		} else {
			lt, oc = deribitLimitTrader(rollToContract, false)
		}
	} else if action == "rollx" {
		// roll from less liquid to more liquid contracts
		if !strings.HasPrefix(contract, "BTC-") {
			log.Fatal("Can only roll from Deribit")
		}

		lt, oc = deribitLimitTrader(contract, true)

		if rollToContract == "BTCUSD" {
			mt = bybitMarketTrader(false)
		} else {
			mt = deribitMarketTrader(rollToContract, false)
		}
	}

	if err := lt(usd, processCb(mt)); err != nil {
		log.Fatal(err)
	}

	log.Debug("Waiting on callbacks")
	wg.Wait()

	log.Info("Done")
}
