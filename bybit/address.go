package bybit

import (
	"errors"
	"net/url"
)

func (b *Bybit) GetAddress() (string, error) {
	var response addressResponse

	err := b.get("/v5/asset/deposit/query-address",
		url.Values{"coin": {"BTC"}}, &response)

	if err != nil {
		return "", err
	}

	for _, chain := range response.Result.Chains {
		if chain.ChainType == "BTC" {
			return chain.AddressDeposit, nil
		}
	}

	return "", errors.New("Address not found")
}
