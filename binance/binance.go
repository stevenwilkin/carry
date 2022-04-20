package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
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

	return body, nil
}

func (b *Binance) marketOrder(usdt float64, buy bool) error {
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
	return b.marketOrder(usdt, true)
}

func (b *Binance) Sell(usdt float64) error {
	return b.marketOrder(usdt, false)
}
