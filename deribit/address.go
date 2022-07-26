package deribit

import (
	"net/url"
)

func (d *Deribit) GetAddress() (string, error) {
	var response addressResponse
	err := d.get(
		"/api/v2/private/get_current_deposit_address",
		url.Values{"currency": {"BTC"}},
		&response)

	if err != nil {
		return "", err
	}

	return response.Result.Address, nil
}
