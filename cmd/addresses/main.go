package main

import (
	"fmt"
	"os"

	"github.com/stevenwilkin/carry/binance"
	"github.com/stevenwilkin/carry/bybit"
	"github.com/stevenwilkin/carry/deribit"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	b := &binance.Binance{
		ApiKey:    os.Getenv("BINANCE_API_KEY"),
		ApiSecret: os.Getenv("BINANCE_API_SECRET")}

	d := &deribit.Deribit{
		ApiId:     os.Getenv("DERIBIT_API_ID"),
		ApiSecret: os.Getenv("DERIBIT_API_SECRET")}

	by := &bybit.Bybit{
		ApiKey:    os.Getenv("BYBIT_API_KEY"),
		ApiSecret: os.Getenv("BYBIT_API_SECRET")}

	btc, _ := b.GetAddress("BTC")
	usdt, _ := b.GetAddress("USDT")

	if btc != "" && usdt != "" {
		fmt.Println("Binance")
		fmt.Printf("        BTC: %s\n", btc)
		fmt.Printf("        USDT: %s\n", usdt)
	}

	btc, _ = d.GetAddress()
	if btc != "" {
		fmt.Println("Deribit")
		fmt.Printf("        BTC: %s\n", btc)
	}

	btc, _ = by.GetAddress()
	if btc != "" {
		fmt.Println("Bybit  ")
		fmt.Printf("        BTC: %s\n", btc)
	}
}
