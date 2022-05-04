package deribit

import (
	"errors"
	"net/url"
	"strconv"
)

func (d *Deribit) MarketOrder(instrument string, amount int, buy, reduce bool) error {
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

	var response orderResponse
	if err := d.get(path, v, &response); err != nil {
		return err
	}

	if response.Error.Message != "" {
		return errors.New(response.Error.Message)
	}

	return nil
}
