package deribit

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Deribit struct {
	ApiId     string
	ApiSecret string
	Test      bool
	Bid       float64
	Ask       float64
}

func (d *Deribit) hostname() string {
	if d.Test {
		return "test.deribit.com"
	} else {
		return "www.deribit.com"
	}
}

func (d *Deribit) get(path string, params url.Values, result interface{}) error {
	u := url.URL{
		Scheme:   "https",
		Host:     d.hostname(),
		Path:     path,
		RawQuery: params.Encode()}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", d.ApiId, d.ApiSecret)))

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Warn(err.Error())
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	json.Unmarshal(body, result)

	return nil
}

func (d *Deribit) subscribe(channels []string) (*websocket.Conn, error) {
	socketUrl := url.URL{Scheme: "wss", Host: d.hostname(), Path: "/ws/api/v2"}

	c, _, err := websocket.DefaultDialer.Dial(socketUrl.String(), nil)
	if err != nil {
		return &websocket.Conn{}, err
	}

	authRequest := requestMessage{
		Method: "/public/auth",
		Params: map[string]interface{}{
			"client_id":     d.ApiId,
			"client_secret": d.ApiSecret,
			"grant_type":    "client_credentials"}}

	if err = c.WriteJSON(authRequest); err != nil {
		return &websocket.Conn{}, err
	}

	request := requestMessage{
		Method: "/private/subscribe",
		Params: map[string]interface{}{
			"channels": channels}}

	if err = c.WriteJSON(request); err != nil {
		return &websocket.Conn{}, err
	}

	ticker := time.NewTicker(10 * time.Second)
	testMessage := requestMessage{Method: "/public/test"}

	go func() {
		for {
			if err = c.WriteJSON(testMessage); err != nil {
				log.WithField("venue", "deribit").Debug("Heartbeat stopping")
				return
			}
			<-ticker.C
		}
	}()

	return c, nil
}

func NewDeribitFromEnv() *Deribit {
	return &Deribit{
		ApiId:     os.Getenv("DERIBIT_API_ID"),
		ApiSecret: os.Getenv("DERIBIT_API_SECRET")}
}
