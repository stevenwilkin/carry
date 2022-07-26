package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/stevenwilkin/carry/binance"
	"github.com/stevenwilkin/carry/bybit"
	"github.com/stevenwilkin/carry/deribit"

	_ "github.com/joho/godotenv/autoload"
)

var (
	wg          sync.WaitGroup
	bBtc, bUsdt string
	dBtc        string
	byBtc       string
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

	wg.Add(4)

	go func() {
		bBtc, _ = b.GetAddress("BTC")
		wg.Done()
	}()

	go func() {
		bUsdt, _ = b.GetAddress("USDT")
		wg.Done()
	}()

	go func() {
		dBtc, _ = d.GetAddress()
		wg.Done()
	}()

	go func() {
		byBtc, _ = by.GetAddress()
		wg.Done()
	}()

	wg.Wait()

	if bBtc != "" || bUsdt != "" {
		fmt.Println("Binance")
		if bBtc != "" {
			fmt.Printf("        BTC: %s\n", bBtc)
		}
		if bUsdt != "" {
			fmt.Printf("        USDT: %s\n", bUsdt)
		}
	}

	if dBtc != "" {
		fmt.Println("Deribit")
		fmt.Printf("        BTC: %s\n", dBtc)
	}

	if byBtc != "" {
		fmt.Println("Bybit  ")
		fmt.Printf("        BTC: %s\n", byBtc)
	}
}
