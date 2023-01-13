package main

import (
	"math"
	"os"
	"regexp"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func testnet() bool {
	return os.Getenv("TESTNET") != ""
}

func validContract(s string) bool {
	// BTCUSD        - Bybit
	// BTCUSD_PERP   - Binance
	// BTC-PERPETUAL - Deribit
	matched, _ := regexp.MatchString("^BTC(USD|USD_PERP|-PERPETUAL|-[0-9]{1,2}[A-Z]{3}[0-9]{2})$", s)
	return matched
}

func initParams() {
	if level, err := log.ParseLevel(os.Getenv("LOG_LEVEL")); err == nil {
		log.SetLevel(level)
	}

	if testnet() {
		log.Warn("Running against testnet")
	}

	if len(os.Args) < 4 {
		log.Fatalf("Usage: %s up|down|roll PARAMS", os.Args[0])
	}

	params = os.Args[1:]

	action, params = params[0], params[1:]
	if matched, _ := regexp.MatchString("^(up|down|roll)$", action); !matched {
		log.Fatalf("Invalid action: %s", action)
	}

	contract, params = params[0], params[1:]
	if !validContract(contract) {
		log.Fatalf("Invalid contract: %s", contract)
	}

	if action == "roll" {
		rollToContract, params = params[0], params[1:]
		if !validContract(rollToContract) {
			log.Fatalf("Invalid contract: %s", rollToContract)
		}

		if contract == rollToContract {
			log.Fatalf("Cannot roll to the same contract")
		}
	}

	if len(params) > 0 {
		var usdStr string
		usdStr, params = params[0], params[1:]
		usdFloat, _ := strconv.ParseFloat(usdStr, 64)
		usd = int(math.Abs(usdFloat))
	}

	if usd == 0 {
		log.Fatal("Number of contracts cannot be zero")
	}

	rounds = 1
	if len(params) > 0 {
		roundsFloat, _ := strconv.ParseFloat(params[0], 64)
		rounds = int(math.Abs(roundsFloat))
	}
}
