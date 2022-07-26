package main

import (
	"fmt"
	"os"

	"github.com/stevenwilkin/carry/binance"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	b := &binance.Binance{
		ApiKey:    os.Getenv("BINANCE_API_KEY"),
		ApiSecret: os.Getenv("BINANCE_API_SECRET")}

	btc, _ := b.GetAddress("BTC")
	usdt, _ := b.GetAddress("USDT")

	if btc != "" && usdt != "" {
		fmt.Println("Binance")
		fmt.Printf("        BTC: %s\n", btc)
		fmt.Printf("        USDT: %s\n", usdt)
	}
}
