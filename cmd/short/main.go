package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/stevenwilkin/carry/deribit"

	_ "github.com/joho/godotenv/autoload"
)

func validContract(s string) bool {
	matched, _ := regexp.MatchString("^BTC-(PERPETUAL|[0-9]{1,2}[A-Z]{3}[0-9]{2})$", s)
	return matched
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s CONTRACT QUANTITY\n", os.Args[0])
		os.Exit(1)
	}

	contract := os.Args[1]
	if !validContract(contract) {
		fmt.Printf("Invalid contract: %s\n", contract)
		os.Exit(1)
	}

	amount, _ := strconv.Atoi(os.Args[2])
	if amount == 0 {
		fmt.Println("Amount cannot be zero")
		os.Exit(1)
	}

	d := deribit.NewDeribitFromEnv()
	if err := d.MarketOrder(contract, amount, false, false); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
