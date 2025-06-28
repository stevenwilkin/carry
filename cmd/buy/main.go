package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/stevenwilkin/carry/binance"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	b := &binance.Binance{
		ApiKey:    os.Getenv("BINANCE_API_KEY"),
		ApiSecret: os.Getenv("BINANCE_API_SECRET")}

	var usdt float64

	if len(os.Args) >= 2 {
		usdt, _ = strconv.ParseFloat(os.Args[1], 64)
	}

	if usdt < 5.0 {
		fmt.Println("Invalid args")
		return
	}

	fmt.Printf("USDT: %.2f\n", usdt)

	if err := b.Buy(usdt); err != nil {
		fmt.Println(err)
	}
}
