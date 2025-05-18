package bybit

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Bybit struct {
	ApiKey    string
	ApiSecret string
	Testnet   bool
	Bid       float64
	Ask       float64
	o         sync.Once
}

func (b *Bybit) hostname() string {
	if b.Testnet {
		return "api-testnet.bybit.com"
	} else {
		return "api.bybit.com"
	}
}

func (b *Bybit) websocketHostname() string {
	if b.Testnet {
		return "stream-testnet.bybit.com"
	} else {
		return "stream.bybit.com"
	}
}

func (b *Bybit) timestamp() string {
	return strconv.FormatInt((time.Now().UnixNano() / int64(time.Millisecond)), 10)
}

func (b *Bybit) get(path string, params url.Values, result interface{}) error {
	query := params.Encode()
	timestamp := b.timestamp()

	u := url.URL{
		Scheme:   "https",
		Host:     b.hostname(),
		Path:     path,
		RawQuery: query}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}

	h := hmac.New(sha256.New, []byte(b.ApiSecret))
	io.WriteString(h, timestamp+b.ApiKey+query)
	signature := fmt.Sprintf("%x", h.Sum(nil))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-BAPI-API-KEY", b.ApiKey)
	req.Header.Set("X-BAPI-SIGN", signature)
	req.Header.Set("X-BAPI-TIMESTAMP", timestamp)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(body, result)

	return nil
}

func (b *Bybit) post(path string, params map[string]interface{}, result interface{}) error {
	jsonRequest, err := json.Marshal(params)
	if err != nil {
		return err
	}

	u := url.URL{
		Scheme: "https",
		Host:   b.hostname(),
		Path:   path}

	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(jsonRequest))
	if err != nil {
		return err
	}

	timestamp := b.timestamp()

	h := hmac.New(sha256.New, []byte(b.ApiSecret))
	io.WriteString(h, timestamp+b.ApiKey+string(jsonRequest))
	signature := fmt.Sprintf("%x", h.Sum(nil))

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-BAPI-API-KEY", b.ApiKey)
	req.Header.Set("X-BAPI-SIGN", signature)
	req.Header.Set("X-BAPI-TIMESTAMP", timestamp)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(body, result)

	return nil
}

func (b *Bybit) subscribe(channels []string) (*websocket.Conn, error) {
	expires := (time.Now().UnixNano() / int64(time.Millisecond)) + 10000

	signatureInput := fmt.Sprintf("GET/realtime%d", expires)
	h := hmac.New(sha256.New, []byte(b.ApiSecret))
	io.WriteString(h, signatureInput)
	signature := fmt.Sprintf("%x", h.Sum(nil))

	v := url.Values{}
	v.Set("api_key", b.ApiKey)
	v.Set("expires", strconv.FormatInt(expires, 10))
	v.Set("signature", signature)

	u := url.URL{
		Scheme:   "wss",
		Host:     b.websocketHostname(),
		Path:     "/realtime",
		RawQuery: v.Encode()}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return &websocket.Conn{}, err
	}

	command := wsCommand{Op: "subscribe", Args: channels}
	if err = c.WriteJSON(command); err != nil {
		return &websocket.Conn{}, err
	}

	return c, nil
}

func NewBybitFromEnv() *Bybit {
	return &Bybit{
		ApiKey:    os.Getenv("BYBIT_API_KEY"),
		ApiSecret: os.Getenv("BYBIT_API_SECRET"),
		Testnet:   os.Getenv("TESTNET") != ""}

}
