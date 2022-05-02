package bybit

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type Bybit struct {
	ApiKey    string
	ApiSecret string
	Testnet   bool
	Bid       float64
	Ask       float64
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

func (b *Bybit) signedUrl(path string, addParams map[string]string) string {
	params := map[string]interface{}{
		"api_key":   b.ApiKey,
		"timestamp": b.timestamp()}

	for k, v := range addParams {
		params[k] = v
	}

	keys := make([]string, len(params))
	i := 0
	query := ""
	for k, _ := range params {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		query += fmt.Sprintf("%s=%v&", k, params[k])
	}
	query = query[0 : len(query)-1]
	h := hmac.New(sha256.New, []byte(b.ApiSecret))
	io.WriteString(h, query)
	query += fmt.Sprintf("&sign=%x", h.Sum(nil))

	u := url.URL{
		Scheme:   "https",
		Host:     b.hostname(),
		Path:     path,
		RawQuery: query}

	return u.String()
}

func (b *Bybit) get(path string, params map[string]string, result interface{}) error {
	u := b.signedUrl(path, params)

	resp, err := http.Get(u)
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
