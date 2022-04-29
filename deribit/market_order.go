package deribit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func (d *Deribit) MarketOrder(instrument string, amount int, buy, reduce bool) (string, error) {
	v := url.Values{}
	v.Set("instrument_name", instrument)
	v.Set("amount", strconv.Itoa(amount))
	v.Set("type", "market")

	path := "/api/v2/private/buy"
	if !buy {
		path = "/api/v2/private/sell"
	}

	if reduce {
		v.Set("reduce_only", "true")
	}

	u := url.URL{
		Scheme:   "https",
		Host:     d.hostname(),
		Path:     path,
		RawQuery: v.Encode()}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", err
	}

	accessToken, err := d.accessToken()
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response orderResponse
	json.Unmarshal(body, &response)

	return response.Result.Order.OrderId, nil
}
