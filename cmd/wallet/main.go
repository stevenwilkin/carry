package main

import (
	"fmt"
	"os"

	"github.com/stevenwilkin/carry/binance"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	b := binance.NewBinanceFromEnv()

	balances, err := b.GetBalances()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, balance := range balances {
		if balance.Asset == "USDT" {
			fmt.Printf("%5s: %-11.2f %11.2f USD\n", balance.Asset, balance.Balance, balance.Value)
		} else {
			fmt.Printf("%5s: %-11.8f %11.2f USD\n", balance.Asset, balance.Balance, balance.Value)
		}
	}
}
