package deribit

import (
	"errors"
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func (d *Deribit) LimitOrder(instrument string, amount int, price float64, buy, reduce bool) (string, error) {
	log.WithFields(log.Fields{
		"venue":      "deribit",
		"instrument": instrument,
		"amount":     amount,
		"price":      price,
	}).Debug("Placing order")

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

	if response.Error.Message != "" {
		return "", errors.New(response.Error.Message)
	}

	return response.Result.Order.OrderId, nil
}

func (d *Deribit) EditOrder(orderId string, amount int, price float64, reduce bool) error {
	log.WithFields(log.Fields{
		"venue":  "deribit",
		"order":  orderId,
		"amount": amount,
		"price":  price,
	}).Debug("Updating order")

	v := url.Values{}
	v.Set("order_id", orderId)
	v.Set("amount", strconv.Itoa(amount))
	v.Set("price", strconv.FormatFloat(price, 'f', 2, 64))
	v.Set("post_only", "true")
	v.Set("reject_post_only", "true")

	if reduce {
		v.Set("reduce_only", "true")
	}

	var response orderResponse
	if err := d.get("/api/v2/private/edit", v, &response); err != nil {
		return err
	}

	if response.Error.Message != "" {
		return errors.New(response.Error.Message)
	}

	return nil
}
