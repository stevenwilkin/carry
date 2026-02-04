package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Binance struct {
	ApiKey    string
	ApiSecret string
	Testnet   bool
}

func (b *Binance) hostname() string {
	if b.Testnet {
		return "testnet.binance.vision"
	} else {
		return "api.binance.com"
	}
}

func (b *Binance) sign(s string) string {
	h := hmac.New(sha256.New, []byte(b.ApiSecret))
	io.WriteString(h, string(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (b *Binance) doRequest(method, path string, values url.Values, sign bool) ([]byte, error) {
	var params string

	if sign {
		timestamp := time.Now().UnixNano() / int64(time.Millisecond)
		values.Set("timestamp", fmt.Sprintf("%d", timestamp))
		input := values.Encode()
		params = fmt.Sprintf("%s&signature=%s", input, b.sign(input))
	} else {
		params = values.Encode()
	}

	u := url.URL{
		Scheme:   "https",
		Host:     b.hostname(),
		Path:     path,
		RawQuery: params}

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return []byte{}, err
	}

	req.Header.Set("X-MBX-APIKEY", b.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode != http.StatusOK {
		var response errorResponse
		json.Unmarshal(body, &response)
		return []byte{}, errors.New(response.Msg)
	}

	return body, nil
}

func (b *Binance) MarketOrder(usdt float64, buy bool) error {
	log.WithFields(log.Fields{
		"venue": "binance",
		"usdt":  usdt,
		"buy":   buy,
	}).Debug("Placing market order")

	side := "BUY"
	if !buy {
		side = "SELL"
	}

	v := url.Values{
		"symbol":        {"BTCUSDT"},
		"side":          {side},
		"type":          {"MARKET"},
		"quoteOrderQty": {strconv.FormatFloat(usdt, 'f', 8, 64)}}

	_, err := b.doRequest("POST", "/api/v3/order", v, true)
	return err
}

func (b *Binance) Buy(usdt float64) error {
	return b.MarketOrder(usdt, true)
}

func (b *Binance) Sell(usdt float64) error {
	return b.MarketOrder(usdt, false)
}

func (b *Binance) GetBalance() (float64, error) {
	body, err := b.doRequest("GET", "/api/v3/account", url.Values{}, true)
	if err != nil {
		return 0, err
	}

	var response accountResponse
	json.Unmarshal(body, &response)

	var usdt float64

	for _, asset := range response.Balances {
		switch asset.Asset {
		case "USDT":
			usdt = asset.Total()
		default:
			continue
		}
	}

	return usdt, nil
}

func (b *Binance) GetPrice(symbol string) (float64, error) {
	params := url.Values{"symbol": {symbol}}
	body, err := b.doRequest("GET", "/api/v3/ticker/price", params, false)
	if err != nil {
		return 0, err
	}

	var response priceResponse
	json.Unmarshal(body, &response)

	return response.PriceFloat(), nil
}

func (b *Binance) GetBalances() ([]balance, error) {
	var wg sync.WaitGroup
	var err error
	var usdtUsd float64

	wg.Add(1)
	go func() {
		usdtUsd, err = b.GetPrice("USDTUSD")
		wg.Done()
	}()

	params := url.Values{"omitZeroBalances": {"true"}}
	body, errGet := b.doRequest("GET", "/api/v3/account", params, true)
	if errGet != nil {
		return []balance{}, err
	}

	var response accountResponse
	json.Unmarshal(body, &response)

	wg.Wait()
	if err != nil {
		return []balance{}, err
	}

	results := []balance{}
	i := 0

	for _, asset := range response.Balances {
		if asset.Asset == "EDG" {
			continue
		}

		results = append(results, balance{
			Asset: asset.Asset, Balance: asset.Total()})

		if asset.Asset == "USDT" {
			results[i].Value = results[i].Balance * usdtUsd
			i++
			continue
		}

		wg.Add(1)
		go func(symbol string, idx int) {
			price, errPrice := b.GetPrice(symbol)
			if err != nil {
				err = errPrice
			} else {
				results[idx].Value = results[idx].Balance * price * usdtUsd
			}
			wg.Done()
		}(asset.Asset+"USDT", i)

		i++
	}

	wg.Wait()
	if err != nil {
		return []balance{}, err
	}

	return results, nil
}

func (b *Binance) GetAddress(coin string) (string, error) {
	network := "ETH"
	if coin == "BTC" {
		network = "BTC"
	}
	params := url.Values{"coin": {coin}, "network": {network}}

	body, err := b.doRequest(
		"GET", "/sapi/v1/capital/deposit/address", params, true)
	if err != nil {
		return "", err
	}

	var response addressResponse
	json.Unmarshal(body, &response)

	return response.Address, nil
}

func NewBinanceFromEnv() *Binance {
	return &Binance{
		ApiKey:    os.Getenv("BINANCE_API_KEY"),
		ApiSecret: os.Getenv("BINANCE_API_SECRET"),
		Testnet:   os.Getenv("TESTNET") != ""}
}
