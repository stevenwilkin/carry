package bybit

func (b *Bybit) GetSize() int {
	var response positionResponse

	err := b.get("/v2/private/position/list",
		map[string]string{"symbol": "BTCUSD"}, &response)

	if err != nil {
		return 0
	}

	return response.Result.Size
}
