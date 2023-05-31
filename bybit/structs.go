package bybit

import (
	"encoding/json"
)

type positionResponse struct {
	Result struct {
		List []struct {
			Size string `json:"size"`
		} `json:"list"`
	} `json:"result"`
}

type orderResponse struct {
	RetCode int    `json:"ret_code"`
	RetMsg  string `json:"ret_msg"`
	Result  struct {
		OrderId string `json:"order_id"`
	} `json:"result"`
}

type addressResponse struct {
	Result struct {
		Chains []struct {
			ChainType      string `json:"chainType"`
			AddressDeposit string `json:"addressDeposit"`
		} `json:"chains"`
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

type orderTopicData struct {
	Topic string `json:"topic"`
	Data  []struct {
		OrderId     string `json:"order_id"`
		OrderStatus string `json:"order_status"`
		Qty         int    `json:"qty"`
		CumExecQty  int    `json:"cum_exec_qty"`
	} `json:"data"`
}
