package deribit

import (
	"net/url"
	"sort"
	"time"
)

func (d *Deribit) GetPositions() []Position {
	var response positionsResponse

	err := d.get("/api/v2/private/get_positions",
		url.Values{"currency": {"BTC"}, "kind": {"future"}}, &response)

	if err != nil {
		return []Position{}
	}

	result := make([]Position, len(response.Result))
	timestamps := make([]int, len(response.Result))
	unsorted := map[int]Position{}

	for i, position := range response.Result {
		t, err := time.Parse("2Jan06", position.InstrumentName[4:])
		if err != nil {
			continue
		}

		timestamp := int(t.Unix())
		timestamps[i] = timestamp
		unsorted[timestamp] = position
	}

	sort.Ints(timestamps)

	for i, timestamp := range timestamps {
		result[i] = unsorted[timestamp]
	}

	return result
}
