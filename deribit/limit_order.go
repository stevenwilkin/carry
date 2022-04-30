package deribit

import (
	"net/url"
	"strconv"
)

func (d *Deribit) LimitOrder(instrument string, amount int, price float64, buy, reduce bool) (string, error) {
	v := url.Values{}
	v.Set("instrument_name", instrument)
	v.Set("amount", strconv.Itoa(amount))
	v.Set("price", strconv.FormatFloat(price, 'f', 2, 64))
	v.Set("post_only", "true")
	v.Set("reject_post_only", "true")

	path := "/api/v2/private/buy"
	if !buy {
		path = "/api/v2/private/sell"
	}

	if reduce {
		v.Set("reduce_only", "true")
	}

	var response orderResponse

	if err := d.get(path, v, &response); err != nil {
		return "", err
	}

	return response.Result.Order.OrderId, nil
}
