package binance_futures

import (
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
)

type BinanceFutures struct {
	ApiKey    string
	ApiSecret string
	Testnet   bool
	Bid       float64
	Ask       float64
}

func (b *BinanceFutures) wsHostname() string {
	if b.Testnet {
		return "dstream.binancefuture.com"
	} else {
		return "dstream.binance.com"
	}
}

func (b *BinanceFutures) subscribe(stream string) (*websocket.Conn, error) {
	u := url.URL{
		Scheme: "wss",
		Host:   b.wsHostname(),
		Path:   fmt.Sprintf("/ws/%s", stream)}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return &websocket.Conn{}, err
	}

	return c, nil
}
