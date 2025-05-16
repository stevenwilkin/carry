package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/stevenwilkin/carry/binance"
	"github.com/stevenwilkin/carry/bybit"
	"github.com/stevenwilkin/carry/deribit"

	_ "github.com/joho/godotenv/autoload"
)

var (
	b  *binance.Binance
	by *bybit.Bybit
	d  *deribit.Deribit

	usdt    float64
	btcusd  int
	futures []deribit.Position
)

func width(s string, x int) string {
	return fmt.Sprintf(fmt.Sprintf("%%-%ds", x), s)
}

func display() {
	fmt.Println("\033[2J\033[H\033[?25l") // clear screen, move cursor to top of screen, hide cursor

	total := usdt + float64(btcusd)

	w := len("USDT:")
	if len(futures) > 0 {
		w = len("BTC-PERPETUAL:")
	} else if btcusd != 0 {
		w = len("BTCUSD:")
	}

	if usdt != 0 {
		fmt.Printf("  %s %6.0f\n", width("USDT:", w), usdt)
	}

	if btcusd != 0 {
		fmt.Printf("  %s %6d\n", width("BTCUSD:", w), btcusd)
	}

	for _, position := range futures {
		total += math.Abs(position.Size)
		fmt.Printf("  %s %6.0f\n",
			width(position.InstrumentName+":", w), math.Abs(position.Size))
	}

	fmt.Printf("  %s %6.0f\n", width("", w), total)
}

func main() {
	b = &binance.Binance{
		ApiKey:    os.Getenv("BINANCE_API_KEY"),
		ApiSecret: os.Getenv("BINANCE_API_SECRET")}

	by = &bybit.Bybit{
		ApiKey:    os.Getenv("BYBIT_API_KEY"),
		ApiSecret: os.Getenv("BYBIT_API_SECRET")}

	d = &deribit.Deribit{
		ApiId:     os.Getenv("DERIBIT_API_ID"),
		ApiSecret: os.Getenv("DERIBIT_API_SECRET")}

	go func() {
		t := time.NewTicker(1 * time.Second)
		var err error

		for {
			usdt, err = b.GetBalance()
			if err != nil {
				panic(err)
			}

			<-t.C
		}
	}()

	go func() {
		t := time.NewTicker(1 * time.Second)

		for {
			btcusd = by.GetSize()
			<-t.C
		}
	}()

	go func() {
		t := time.NewTicker(1 * time.Second)

		for {
			futures = d.GetPositions()
			<-t.C
		}
	}()

	t := time.NewTicker(100 * time.Millisecond)

	for {
		display()
		<-t.C
	}
}
