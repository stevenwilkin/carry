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
)

type Bybit struct {
	ApiKey    string
	ApiSecret string
	Testnet   bool
}

func (b *Bybit) hostname() string {
	if b.Testnet {
		return "api-testnet.bybit.com"
	} else {
		return "api.bybit.com"
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
