package bybit

import (
	"errors"
)

func (b *Bybit) GetAddress() (string, error) {
	var response addressResponse

	err := b.get("/v5/asset/deposit/query-address",
		map[string]string{"coin": "BTC"}, &response)

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
