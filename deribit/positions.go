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

	timestamps := []int{}
	unsorted := map[int]Position{}

	for _, position := range response.Result {
		if position.Size == 0 {
			continue
		}

		t, err := time.Parse("2Jan06", position.InstrumentName[4:])
		if err != nil {
			continue
		}

		timestamp := int(t.Unix())
		timestamps = append(timestamps, timestamp)
		unsorted[timestamp] = position
	}

	result := make([]Position, len(timestamps))
	sort.Ints(timestamps)

	for i, timestamp := range timestamps {
		result[i] = unsorted[timestamp]
	}

	return result
}
