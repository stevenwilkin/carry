package bybit

import (
	"net/url"
	"strconv"
)

func (b *Bybit) GetSize() int {
	params := url.Values{
		"category": {"inverse"},
		"symbol":   {"BTCUSD"}}

	var response positionResponse
	err := b.get("/v5/position/list", params, &response)

	if err != nil {
		return 0
	}

	if len(response.Result.List) != 1 {
		return 0
	}

	size, _ := strconv.Atoi(response.Result.List[0].Size)
	return size
}
