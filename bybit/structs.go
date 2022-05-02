package bybit

import (
	"encoding/json"
)

type positionResponse struct {
	Result struct {
		Size int `json:"size"`
	} `json:"result"`
}

type orderResponse struct {
	Result struct {
		OrderId string `json:"order_id"`
	} `json:"result"`
}

type wsCommand struct {
	Op   string   `json:"op"`
	Args []string `json:"args"`
}

type wsResponse struct {
	Topic string          `json:"topic"`
	Type  string          `json:"type"`
	Data  json.RawMessage `json:"data"`
}

type order struct {
	Id    int64  `json:"id"`
	Price string `json:"price"`
	Side  string `json:"side"`
}

type snapshotData []order

type updateData struct {
	Delete []order `json:"delete"`
	Insert []order `json:"insert"`
}
