package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/stevenwilkin/carry/binance"
	"github.com/stevenwilkin/carry/bybit"
	"github.com/stevenwilkin/carry/deribit"
	"github.com/stevenwilkin/carry/feed"

	_ "github.com/joho/godotenv/autoload"
)

var (
	h  = feed.NewHandler()
	b  = binance.NewBinanceFromEnv()
	by = bybit.NewBybitFromEnv()
	d  = deribit.NewDeribitFromEnv()

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

func displayFailing() {
	fmt.Println("\033[2J\033[H\033[?25l") // clear screen, move cursor to top of screen, hide cursor
	fmt.Println("Feed failing...")
}

func exitFailed() {
	fmt.Println("Feed failed")
	os.Exit(1)
}

func main() {
	h.Add(feed.NewFeed(feed.Poll(b.GetBalance), feed.SetValue(&usdt)))
	h.Add(feed.NewFeed(feed.Poll(by.GetSize), feed.SetValue(&btcusd)))
	h.Add(feed.NewFeed(feed.Poll(d.GetPositions), feed.SetValue(&futures)))

	t := time.NewTicker(100 * time.Millisecond)

	for {
		if h.Failed() {
			exitFailed()
		} else if h.Failing() {
			displayFailing()
		} else {
			display()
		}

		<-t.C
	}
}
